package chatcontext

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	chatcommands "github.com/xfrr/randomtalk/internal/chat/application/commands"
	chatqueries "github.com/xfrr/randomtalk/internal/chat/application/queries"
	chatconfig "github.com/xfrr/randomtalk/internal/chat/config"
	chathttp "github.com/xfrr/randomtalk/internal/chat/infrastructure/http"
	"github.com/xfrr/randomtalk/internal/shared/env"
	"github.com/xfrr/randomtalk/internal/shared/logging"
	xotel "github.com/xfrr/randomtalk/internal/shared/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrEmptyConfig = errors.New("service config is empty, please provide a valid config")
)

// Service defines the dependencies that can be overridden
// when initializing a new chat service.
type Service struct {
	version          string
	traceProvider    trace.TracerProvider
	config           chatconfig.Config
	logger           *zerolog.Logger
	natsConnection   *nats.Conn
	httpWebsocketHub *chathttp.Hub
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		s.shutdown()
	}()

	go s.httpWebsocketHub.Run()

	// serve http and websocket
	go func() {
		s.logger.Info().
			Str("address", s.config.WebsocketServer.Address).
			Str("path", s.config.WebsocketServer.Path).
			Str("url", fmt.Sprintf("ws://%s%s", s.config.WebsocketServer.Address, s.config.WebsocketServer.Path)).
			Msg("starting websocket server")
		err := chathttp.ServeHTTP(s.config.WebsocketServer, s.httpWebsocketHub.Handle)
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
		service.logger.Fatal().Err(err).Msg("failed to initialize chat service")
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

	if svc.logger == nil {
		logger := logging.NewLogger(
			svc.config.ServiceName,
			env.Environment(svc.config.ServiceEnvironment),
			svc.config.Logging.Level,
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

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      svc.config.NotificationStream.Name,
		Subjects:  []string{"randomtalk.notifications.chat.>"},
		Retention: jetstream.LimitsPolicy,
		MaxAge:    5 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	cmdbus, querybus := chatcommands.InitCommandBus(ctx), chatqueries.InitQueryBus(ctx)
	svc.httpWebsocketHub = chathttp.NewHub(
		cmdbus,
		querybus,
		chathttp.WithLogger(*svc.logger),
	)

	return svc, nil
}

func initOtelTraces(ctx context.Context, config chatconfig.Config, serviceVersion string) (trace.TracerProvider, error) {
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

func (s *Service) setupNatsConnection(config chatconfig.Config) error {
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
