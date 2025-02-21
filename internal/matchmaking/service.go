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
	matchcommands "github.com/xfrr/randomtalk/internal/matchmaking/application/commands"
	matcheventhandlers "github.com/xfrr/randomtalk/internal/matchmaking/application/handlers"
	matchqueries "github.com/xfrr/randomtalk/internal/matchmaking/application/queries"
	matchmakingconfig "github.com/xfrr/randomtalk/internal/matchmaking/config"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	matchgrpc "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/grpc"
	matchmakinginmemory "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/memory"
	matchnats "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/nats"
	matchmakingtrace "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/tracing"
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
	config             matchmakingconfig.Config
	version            string
	logger             *zerolog.Logger
	natsConnection     *nats.Conn
	grpcServer         *matchgrpc.Server
	grpcServerCloser   func()
	matchmakingService matchdomain.MatchmakingProcessor
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
		config:  matchmakingconfig.Config{},
		version: "development",
	}

	for _, opt := range opts {
		opt(svc)
	}

	if svc.logger == nil {
		logger := logging.NewLogger(
			svc.config.ServiceName,
			env.Environment(svc.config.ServiceEnvironment),
			svc.config.Logging.Level,
		)

		svc.logger = &logger
	}

	if svc.config == (matchmakingconfig.Config{}) {
		svc.config = matchmakingconfig.MustLoadFromEnv()
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

	err = matchnats.CreateMatchmakingUserMatchRequestsStream(ctx, js)
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

	cmdbus := matchcommands.InitCommandBus(ctx, svc.matchmakingService)
	querybus := matchqueries.InitQueryBus(ctx)

	if svc.config.GrpcServerEnabled {
		svc.grpcServer, svc.grpcServerCloser, err = matchgrpc.NewServer(cmdbus, querybus)
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
	chatNotificationConsumer, err := s.initChatNotificationConsumer(ctx)
	if err != nil {
		return err
	}

	go s.startChatNotificationConsumer(ctx, s.matchmakingService, chatNotificationConsumer)

	go s.startGrpcServer()

	return nil
}

func (s *Service) startChatNotificationConsumer(
	ctx context.Context,
	matchmakerService matchdomain.MatchmakingProcessor,
	consumer *matchnats.MessagingEventConsumer,
) {
	// create user match request event handler
	userMatchRequestHandler := matcheventhandlers.NewUserMatchRequestedEventHandler(matchmakerService, s.logger)

	s.logger.Debug().
		Str("consumer_name", s.config.ChatNotificationsConsumerConfig.Name).
		Str("stream_name", s.config.ChatNotificationsConsumerConfig.StreamName).
		Msg("subscribing chat notifications")
	if err := messaging.HandleEvents(
		ctx,
		s.logger,
		consumer,
		userMatchRequestHandler.Handle,
	); err != nil {
		s.logger.Error().Err(err).Msg("failed to start chat notification event handler")
	}
}

func (s *Service) initMatchNotificationsChannel(ctx context.Context, js jetstream.JetStream) (matchdomain.NotificationsChannel, error) {
	var matchNotificationsChannel matchdomain.NotificationsChannel
	matchNotificationsChannel, err := matchnats.CreateMatchNotificationsChannel(ctx, js, s.logger)
	if err != nil {
		return nil, err
	}

	// wrap match notifications channel with tracing
	matchNotificationsChannel = matchmakingtrace.WrapMatchNotificationsChannel(
		matchNotificationsChannel,
		s.traceProvider,
	)
	return matchNotificationsChannel, nil
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

func (s *Service) initChatNotificationConsumer(ctx context.Context) (*matchnats.MessagingEventConsumer, error) {
	chatNotificationConsumer, err := matchnats.SetupMessagingEventConsumer(
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

func (s *Service) initUserInMemoryStore(_ context.Context) (matchdomain.UserStore, error) {
	userStore := matchmakinginmemory.NewUserStore(s.logger)
	return matchmakingtrace.WrapUserStore(userStore, s.traceProvider), nil
}

func (s *Service) initMatchmakerService(
	ctx context.Context,
	js jetstream.JetStream,
	matchRepo matchdomain.MatchRepository,
) (matchdomain.MatchmakingProcessor, error) {
	// userStore, err := s.initUserInMemoryStore(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	userStore, err := matchnats.NewUserStore(ctx, js)
	if err != nil {
		return nil, err
	}

	// freeUsersQueue, err := s.initMatchmakingQueue(ctx, js)
	// if err != nil {
	// 	return nil, err
	// }

	matchNotificationsChannel, err := s.initMatchNotificationsChannel(ctx, js)
	if err != nil {
		return nil, err
	}

	stableMatcher := matchdomain.NewGaleShapleyStableMatcher()

	var matchService matchdomain.MatchmakingProcessor
	matchService, err = matchdomain.NewUserMatchProcessor(
		matchRepo,
		userStore,
		stableMatcher,
		matchNotificationsChannel,
		matchdomain.WithLogger(s.logger),
	)

	matchService = matchmakingtrace.WrapMatchmakingService(
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

func initOtelTraces(ctx context.Context, config matchmakingconfig.Config, serviceVersion string) (trace.TracerProvider, error) {
	traceProvider, err := xotel.InitTracerProvider(ctx,
		xotel.WithServiceName(config.ServiceName),
		xotel.WithServiceVersion(serviceVersion),
		xotel.WithServiceEnvironment(env.Environment(config.ServiceEnvironment)),
		xotel.WithEndpointURL(config.OpenTelemetry.CollectorEndpoint),
	)
	if err != nil {
		return nil, err
	}
	return traceProvider, nil
}

func (s *Service) setupNatsConnection(config matchmakingconfig.Config) error {
	var err error
	s.natsConnection, err = nats.Connect(config.Nats.URI,
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

func (s *Service) initMatchRepository(ctx context.Context, js jetstream.JetStream) (matchdomain.MatchRepository, error) {
	var matchRepo matchdomain.MatchRepository
	matchRepo, err := matchnats.NewMatchStreamRepository(
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

	return matchmakingtrace.WrapMatchRepository(matchRepo, s.traceProvider), nil
}
