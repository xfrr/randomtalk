package matchmaking

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	matchCommands "github.com/xfrr/randomtalk/internal/matchmaking/application/commands"
	matchQueries "github.com/xfrr/randomtalk/internal/matchmaking/application/queries"
	matchmakingConfig "github.com/xfrr/randomtalk/internal/matchmaking/config"
	matchmakingDomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	matchmakingGrpc "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/grpc"
	matchmakingHandlers "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/handlers"
	matchmakingInmemory "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/memory"
	matchmakingNats "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/nats"
	matchmakingTrace "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/tracing"
	"github.com/xfrr/randomtalk/internal/shared/env"
	"github.com/xfrr/randomtalk/internal/shared/logging"
	"github.com/xfrr/randomtalk/internal/shared/messaging"
	xotel "github.com/xfrr/randomtalk/internal/shared/otel"
	"github.com/xfrr/randomtalk/internal/shared/xnats"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrEmptyConfig = errors.New("service config is empty, please provide a valid config")
)

// Service defines the dependencies that can be overridden
// when initializing a new matchmaking service.
type Service struct {
	traceProvider      trace.TracerProvider
	config             matchmakingConfig.Config
	version            string
	logger             *zerolog.Logger
	natsConnection     *nats.Conn
	grpcServer         *matchmakingGrpc.Server
	grpcServerCloser   func()
	matchmakingService matchmakingDomain.MatchmakingProcessor
}

// MustInitService initializes a new matchmaking service with the provided options.
// If an error occurs during initialization, the service will log the error and exit.
func MustInitService(ctx context.Context, opts ...InitOption) Service {
	service, err := NewService(opts...)
	if err != nil {
		service.logger.Fatal().Err(err).Msg("failed to initialize matchmaking service")
	}

	service.logger.Info().Str("version", service.version).Msg("starting matchmaking service")

	if startErr := service.start(ctx); startErr != nil {
		service.logger.Fatal().Err(startErr).Msg("failed to start matchmaking service")
	}

	return *service
}

// NewService creates a new matchmaking service with the provided options.
func NewService(opts ...InitOption) (*Service, error) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc := &Service{
		config:  matchmakingConfig.Config{},
		version: "development",
	}

	for _, opt := range opts {
		opt(svc)
	}
	if svc.config == (matchmakingConfig.Config{}) {
		svc.config = matchmakingConfig.MustLoadFromEnv()
	}

	if svc.logger == nil {
		logger := logging.NewLogger(
			svc.config.ServiceName,
			env.Environment(svc.config.ServiceEnvironment),
			svc.config.LoggingConfig.Level,
		)

		svc.logger = &logger
	}

	if svc.traceProvider == nil {
		svc.traceProvider, err = initOtelTraces(ctx, svc.config, svc.version)
		if err != nil {
			return svc, err
		}
	}

	err = svc.setupNatsConnection(svc.config)
	if err != nil {
		return svc, err
	}

	js, err := jetstream.New(svc.natsConnection)
	if err != nil {
		return svc, err
	}

	err = matchmakingNats.CreateMatchmakingUserMatchRequestsStream(ctx, js)
	if err != nil {
		return nil, err
	}

	// TODO: move this to the chat service
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      svc.config.ChatNotificationsConsumerConfig.StreamName,
		Retention: jetstream.LimitsPolicy,
		Subjects:  []string{"randomtalk.notifications.chat.>"},
		MaxAge:    5 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	matchRepository, err := svc.initMatchRepository(ctx, js)
	if err != nil {
		return svc, err
	}

	svc.matchmakingService, err = svc.initMatchmakerService(
		ctx,
		js,
		matchRepository,
	)
	if err != nil {
		return svc, err
	}

	cmdbus := matchCommands.InitCommandBus(ctx, svc.matchmakingService)
	querybus := matchQueries.InitQueryBus(ctx)

	if svc.config.GrpcServerEnabled {
		svc.grpcServer, svc.grpcServerCloser, err = matchmakingGrpc.NewServer(cmdbus, querybus)
		if err != nil {
			return svc, fmt.Errorf("failed to create gRPC server: %w", err)
		}
	}

	return svc, nil
}

// shutdown closes all the resources used by the matchmaking service.
func (s *Service) Shutdown() {
	if s.natsConnection != nil {
		s.natsConnection.Close()
	}
	if s.grpcServerCloser != nil {
		s.grpcServerCloser()
	}
}

func (s *Service) start(ctx context.Context) error {
	go s.startChatNotificationConsumer(ctx, s.matchmakingService)

	go s.startGrpcServer()

	return nil
}

func (s *Service) startChatNotificationConsumer(
	ctx context.Context,
	matchmakerService matchmakingDomain.MatchmakingProcessor,
) {
	consumer, err := s.initChatNotificationConsumer(ctx)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to initialize chat notification consumer")
		return
	}

	// create user match request event handler
	userMatchRequestHandler := matchmakingHandlers.NewUserMatchRequestedEventHandler(matchmakerService, s.logger)

	s.logger.Debug().
		Str("consumer_name", s.config.ChatNotificationsConsumerConfig.Name).
		Str("stream_name", s.config.ChatNotificationsConsumerConfig.StreamName).
		Msg("subscribing chat notifications")
	if err = messaging.HandleEvents(
		ctx,
		s.logger,
		consumer,
		userMatchRequestHandler.Handle,
	); err != nil {
		s.logger.Error().Err(err).Msg("failed to start chat notification event handler")
	}
}

func (s *Service) startGrpcServer() {
	if s.grpcServer == nil {
		return
	}

	s.logger.Info().Int("port", s.config.GrpcServerPort).Msg("starting gRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GrpcServerPort))
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to listen")
	}

	if err = s.grpcServer.Serve(lis); err != nil {
		s.logger.Fatal().Err(err).Msg("failed to serve")
	}

	s.logger.Info().Msg("gRPC server stopped")
}

func (s *Service) initChatNotificationConsumer(ctx context.Context) (*xnats.MessagingEventConsumer, error) {
	chatNotificationConsumer, err := xnats.CreateMessagingEventConsumer(
		ctx,
		s.natsConnection,
		s.logger,
		s.config.ChatNotificationsConsumerConfig.StreamName,
		jetstream.ConsumerConfig{
			Name:           s.config.ChatNotificationsConsumerConfig.Name,
			Durable:        s.config.ChatNotificationsConsumerConfig.Name,
			AckPolicy:      jetstream.AckExplicitPolicy,
			DeliverPolicy:  jetstream.DeliverAllPolicy,
			AckWait:        30 * time.Second, // TODO: Adjust based on environment settings
			MaxDeliver:     3,
			MaxAckPending:  50, // TODO: Adjust based on environment settings
			FilterSubjects: []string{"randomtalk.notifications.chat.>"},
			BackOff: []time.Duration{
				500 * time.Millisecond,
				1 * time.Second,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return chatNotificationConsumer, nil
}

func (s *Service) initUserInMemoryStore(_ context.Context) (matchmakingDomain.UserStore, error) {
	userStore := matchmakingInmemory.NewUserStore(s.logger)
	return matchmakingTrace.WrapUserStore(userStore, s.traceProvider), nil
}

func (s *Service) initMatchmakerService(
	ctx context.Context,
	js jetstream.JetStream,
	matchRepo matchmakingDomain.MatchRepository,
) (matchmakingDomain.MatchmakingProcessor, error) {
	// userStore, err := s.initUserInMemoryStore(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	userStore, err := matchmakingNats.NewUserStore(ctx, js)
	if err != nil {
		return nil, err
	}

	stableMatcher := matchmakingDomain.NewGaleShapleyStableMatcher()

	var matchService matchmakingDomain.MatchmakingProcessor
	matchService, err = matchmakingDomain.NewUserMatchProcessor(
		matchRepo,
		userStore,
		stableMatcher,
		matchmakingDomain.WithLogger(s.logger),
	)

	matchService = matchmakingTrace.WrapMatchmakingService(
		matchService,
		s.traceProvider,
	)

	if err != nil {
		return nil, err
	}
	return matchService, nil
}

// func (s *Service) processMatchRequests(ctx context.Context, matchService matchdomain.MatchmakingService) {
// 	err := matchService.ProcessMatchQueue(ctx)
// 	if err != nil {
// 		s.logger.Error().Err(err).Msg("failed to process match requests")
// 	}
// }

func initOtelTraces(ctx context.Context, config matchmakingConfig.Config, serviceVersion string) (trace.TracerProvider, error) {
	traceProvider, err := xotel.InitTracerProvider(ctx,
		xotel.WithServiceName(config.ServiceName),
		xotel.WithServiceVersion(serviceVersion),
		xotel.WithServiceEnvironment(env.Environment(config.ServiceEnvironment)),
		xotel.WithEndpointURL(config.Observability.OTELCollectorEndpoint),
	)
	if err != nil {
		return nil, err
	}
	return traceProvider, nil
}

func (s *Service) setupNatsConnection(config matchmakingConfig.Config) error {
	var err error
	s.natsConnection, err = nats.Connect(config.NatsConfig.URI,
		nats.ReconnectWait(5*time.Second),
		nats.MaxReconnects(-1),
		nats.PingInterval(10*time.Second),
		nats.MaxPingsOutstanding(3),
		nats.Timeout(5*time.Second),
		nats.NoEcho(),
		nats.ReconnectJitter(50*time.Millisecond, 1*time.Second),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			s.logger.Error().Err(err).Msg("NATS error")
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			s.logger.Info().Msg("reconnected to NATS")
		}),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) initMatchRepository(ctx context.Context, js jetstream.JetStream) (matchmakingDomain.MatchRepository, error) {
	var matchRepo matchmakingDomain.MatchRepository
	matchRepo, err := matchmakingNats.NewMatchStreamRepository(
		ctx, js, xnats.
			NewStreamConfig("randomtalk_matchmaking_match_events", "randomtalk.matchmaking.matches.>").
			WithDenyDelete().
			// WithDenyPurge(), // TODO: Adjust based on environment settings
			WithReplicas(1). // TODO: Adjust based on environment settings
			WithDiscardPolicy(jetstream.DiscardOld).
			WithMaxAge(24*7*time.Hour). // 1 week
			WithMaxBytes(1<<30),        // 1 GB
	)
	if err != nil {
		return nil, err
	}

	return matchmakingTrace.WrapMatchRepository(matchRepo, s.traceProvider), nil
}
