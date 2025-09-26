package matchmaking

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	"github.com/xfrr/randomtalk/internal/shared/env"
	"github.com/xfrr/randomtalk/internal/shared/logging"
	"github.com/xfrr/randomtalk/internal/shared/messaging"
	xnats "github.com/xfrr/randomtalk/internal/shared/nats"

	commands "github.com/xfrr/randomtalk/internal/matchmaking/application/commands"
	config "github.com/xfrr/randomtalk/internal/matchmaking/config"
	domain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	handlers "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/handlers"
	inMemoryAdapter "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/memory"
	natsAdapter "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/nats"
	tracing "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/tracing"
	xotel "github.com/xfrr/randomtalk/internal/shared/otel"
)

var (
	ErrEmptyConfig = errors.New("service config is empty, please provide a valid config")
)

// Service defines the dependencies that can be overridden
// when initializing a new matchmaking service.
type Service struct {
	traceProvider      trace.TracerProvider
	config             config.Config
	version            string
	logger             *zerolog.Logger
	natsConnection     *nats.Conn
	matchmakingService domain.MatchmakingProcessor
	cmdbus             commands.CommandBus
	closers            []func()
}

// MustInitService initializes a new matchmaking service with the provided options.
// If an error occurs during initialization, the service will log the error and exit.
func MustInitService(ctx context.Context, opts ...InitOption) *Service {
	service, err := NewService(opts...)
	if err != nil {
		service.logger.Fatal().Err(err).Msg("failed to initialize matchmaking service")
	}

	service.logger.Info().Str("version", service.version).Msg("starting matchmaking service")
	service.start(ctx)
	return service
}

// NewService creates a new matchmaking service with the provided options.
func NewService(opts ...InitOption) (*Service, error) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc := &Service{
		config:  config.Config{},
		version: "development",
	}

	for _, opt := range opts {
		opt(svc)
	}
	if svc.config == (config.Config{}) {
		svc.config = config.MustLoadFromEnv()
	}

	if svc.logger == nil {
		logger := logging.NewLogger(
			svc.config.ServiceName,
			env.Environment(svc.config.ServiceEnvironment),
			svc.config.LogLevel,
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

	err = natsAdapter.CreateMatchmakingUserMatchRequestsStream(ctx, js)
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

	var cmdbusCloser func()
	svc.cmdbus, cmdbusCloser = commands.InitCommandBus(ctx, svc.matchmakingService)
	svc.closers = append(svc.closers, cmdbusCloser)

	return svc, nil
}

// shutdown closes all the resources used by the matchmaking service.
func (s *Service) Shutdown() {
	for _, closer := range s.closers {
		closer()
	}
}

func (s *Service) start(ctx context.Context) {
	go s.startChatNotificationConsumer(ctx, s.matchmakingService)
}

func (s *Service) startChatNotificationConsumer(ctx context.Context, mp domain.MatchmakingProcessor) {
	consumer, err := s.initChatNotificationConsumer(ctx)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to initialize chat notification consumer")
		return
	}

	// create user match request event handler
	userMatchRequestHandler := handlers.NewUserMatchRequestedEventHandler(mp, s.logger)

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
			FilterSubjects: []string{"randomtalk.chat.notifications.>"},
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

func (s *Service) initUserInMemoryStore(_ context.Context) (domain.UserStore, error) {
	userStore := inMemoryAdapter.NewUserStore(s.logger)
	return tracing.WrapUserStore(userStore, s.traceProvider), nil
}

func (s *Service) initMatchmakerService(
	ctx context.Context,
	js jetstream.JetStream,
	matchRepo domain.MatchRepository,
) (domain.MatchmakingProcessor, error) {
	userStore, err := natsAdapter.NewUserStore(ctx, js)
	if err != nil {
		return nil, err
	}

	stableMatcher := domain.NewGaleShapleyStableMatcher()

	var matchService domain.MatchmakingProcessor
	matchService, err = domain.NewUserMatchProcessor(
		matchRepo,
		userStore,
		stableMatcher,
		domain.WithLogger(s.logger),
	)

	matchService = tracing.WrapMatchmakingService(
		matchService,
		s.traceProvider,
	)

	if err != nil {
		return nil, err
	}
	return matchService, nil
}

func initOtelTraces(ctx context.Context, config config.Config, serviceVersion string) (trace.TracerProvider, error) {
	traceProvider, err := xotel.InitTracerProvider(ctx,
		xotel.WithServiceName(config.ServiceName),
		xotel.WithServiceVersion(serviceVersion),
		xotel.WithServiceEnvironment(env.Environment(config.ServiceEnvironment)),
		xotel.WithEndpointURL(config.OTELCollectorEndpoint),
	)
	if err != nil {
		return nil, err
	}
	return traceProvider, nil
}

func (s *Service) setupNatsConnection(config config.Config) error {
	var err error
	s.natsConnection, err = nats.Connect(config.NatsURI,
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
		nats.ReconnectHandler(func(_ *nats.Conn) {
			s.logger.Info().Msg("reconnected to NATS")
		}),
	)
	if err != nil {
		return err
	}

	s.registerCloser(func() {
		if s.natsConnection != nil && !s.natsConnection.IsClosed() {
			s.natsConnection.Close()
		}
	})

	s.logger.Info().
		Str("nats_url", config.NatsConfig.NatsURI).
		Msg("connected to NATS server")
	return nil
}

func (s *Service) initMatchRepository(ctx context.Context, js jetstream.JetStream) (domain.MatchRepository, error) {
	var matchRepo domain.MatchRepository
	matchRepo, err := natsAdapter.NewMatchStreamRepository(
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

	return tracing.WrapMatchRepository(matchRepo, s.traceProvider), nil
}

func (s *Service) registerCloser(closer func()) {
	if closer == nil {
		s.closers = make([]func(), 0)
	}

	s.closers = append(s.closers, closer)
}
