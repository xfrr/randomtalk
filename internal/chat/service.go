package chatcontext

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/xfrr/randomtalk/internal/shared/env"
	"github.com/xfrr/randomtalk/internal/shared/logging"
	"go.opentelemetry.io/otel/trace"

	chatcommands "github.com/xfrr/randomtalk/internal/chat/application/commands"
	chatqueries "github.com/xfrr/randomtalk/internal/chat/application/queries"
	chatconfig "github.com/xfrr/randomtalk/internal/chat/config"
	chathttp "github.com/xfrr/randomtalk/internal/chat/infrastructure/http"
	chatnats "github.com/xfrr/randomtalk/internal/chat/infrastructure/nats"
	xnats "github.com/xfrr/randomtalk/internal/shared/nats"
	xotel "github.com/xfrr/randomtalk/internal/shared/otel"
)

var (
	ErrEmptyConfig = errors.New("service config is empty, please provide a valid config")
)

// Service defines the dependencies that can be overridden
// when initializing a new chat service.
type Service struct {
	version                    string
	traceProvider              trace.TracerProvider
	config                     chatconfig.Config
	logger                     *zerolog.Logger
	natsConnection             *nats.Conn
	cmdbus                     chatcommands.CommandBus
	querybus                   chatqueries.QueryBus
	matchNotificationsConsumer *xnats.MessagingEventConsumer
	httpWebsocketHub           *chathttp.Hub
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		s.shutdown()
	}()

	go s.httpWebsocketHub.Run(ctx)

	// serve http and websocket
	go func() {
		s.logger.Info().
			Str("address", s.config.HubWebsocketServer.Address).
			Str("path", s.config.HubWebsocketServer.Path).
			Str("url", fmt.Sprintf("ws://%s%s", s.config.HubWebsocketServer.Address, s.config.HubWebsocketServer.Path)).
			Msg("starting websocket server")
		err := chathttp.ServeHTTP(s.config.HubWebsocketServer, s.httpWebsocketHub.Handle)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to start http server")
		}
	}()
}

func (s *Service) shutdown() {
	if s.natsConnection != nil {
		s.natsConnection.Close()
	}
}

// MustInitService initializes a new chat service with the provided options.
// If an error occurs during initialization, the service will log the error and exit.
func MustInitService(opts ...InitOption) Service {
	service, err := NewService(opts...)
	if err != nil {
		service.logger.Fatal().
			Any("config", service.config).
			Err(err).
			Msg("failed to initialize chat service")
	}
	return *service
}

// NewService creates a new chat service with the provided options.
func NewService(opts ...InitOption) (*Service, error) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	svc := &Service{
		config:  chatconfig.Config{},
		version: "development",
	}

	for _, opt := range opts {
		opt(svc)
	}

	if svc.config == (chatconfig.Config{}) {
		svc.config = chatconfig.MustLoadFromEnv()
	}

	if svc.logger == nil {
		logger := logging.NewLogger(
			svc.config.ServiceName,
			env.Environment(svc.config.ServiceEnvironment),
			svc.config.LoggingConfig.Level,
		)

		svc.logger = &logger
	}

	if svc.config == (chatconfig.Config{}) {
		svc.config = chatconfig.MustLoadFromEnv()
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

	chatSessionRepo, err := chatnats.NewChatSessionRepository(
		ctx,
		js,
		xnats.
			NewStreamConfig(svc.config.ChatSessionStreamConfig.Name, "randomtalk.chat.sessions.>").
			WithReplicas(1).
			WithMaxAge(24*time.Hour).
			WithRetention(jetstream.LimitsPolicy),
	)
	if err != nil {
		return nil, err
	}
	//
	// TODO: Create notification stream and subscribe for chat session domain events
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      svc.config.NotificationStreamConfig.Name,
		Subjects:  []string{"randomtalk.chat.notifications.>"},
		Retention: jetstream.LimitsPolicy,
		MaxAge:    5 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	matchRequester := chatnats.NewMatchRequester(svc.config.NotificationStreamConfig.Name, js)

	svc.cmdbus, svc.querybus = chatcommands.InitCommandBus(ctx, chatSessionRepo, matchRequester, *svc.logger), chatqueries.InitQueryBus(ctx)

	matchNotificationsConsumer, err := svc.initMatchNotificationsConsumer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize match notifications consumer: %w", err)
	}

	svc.httpWebsocketHub = chathttp.NewHub(
		svc.cmdbus,
		svc.querybus,
		matchNotificationsConsumer,
		chathttp.WithLogger(*svc.logger),
	)

	return svc, nil
}

func initOtelTraces(ctx context.Context, config chatconfig.Config, serviceVersion string) (trace.TracerProvider, error) {
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

func (s *Service) setupNatsConnection(config chatconfig.Config) error {
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

func (s *Service) initMatchNotificationsConsumer(ctx context.Context) (*xnats.MessagingEventConsumer, error) {
	chatNotificationConsumer, err := xnats.CreateMessagingEventConsumer(
		ctx,
		s.natsConnection,
		s.logger,
		s.config.MatchNotificationsConsumerConfig.StreamName,
		jetstream.ConsumerConfig{
			Name:           s.config.MatchNotificationsConsumerConfig.Name,
			Durable:        s.config.MatchNotificationsConsumerConfig.Name,
			AckPolicy:      jetstream.AckExplicitPolicy,
			DeliverPolicy:  jetstream.DeliverAllPolicy,
			AckWait:        15 * time.Second, // TODO: Adjust based on environment settings
			MaxDeliver:     3,
			MaxAckPending:  50, // TODO: Adjust based on environment settings
			FilterSubjects: []string{"randomtalk.matchmaking.matches.>"},
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

// func (s *Service) startMatchNotificationsConsumer(
// 	ctx context.Context,
// ) {
// 	// create user match request event handler
// 	matchNotificationsHandler := chatinfrahandlers.NewUserMatchedNotificationHandler(s.cmdbus, s.logger)

// 	s.logger.Debug().
// 		Str("consumer_name", s.config.MatchNotificationsConsumerConfig.Name).
// 		Str("stream_name", s.config.MatchNotificationsConsumerConfig.StreamName).
// 		Msg("subscribing to match notifications...")
// 	if err := messaging.HandleEvents(
// 		ctx,
// 		s.logger,
// 		s.matchNotificationsConsumer,
// 		matchNotificationsHandler.Handle,
// 	); err != nil {
// 		s.logger.Error().Err(err).Msg("failed to start match notifications consumer")
// 	}
// }
