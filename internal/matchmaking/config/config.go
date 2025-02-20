package matchmakingconfig

import "github.com/kelseyhightower/envconfig"

const envPrefix = "RANDOMTALK_MATCHMAKING"

type Config struct {
	ServiceName        string `envconfig:"SERVICE_NAME" default:"randomtalk-matchmaking"`
	ServiceEnvironment string `envconfig:"SERVICE_ENVIRONMENT" default:"development"`

	Persistence                     Persistence                     `envconfig:"-"` // This is not loaded from the environment.
	OpenTelemetry                   OpenTelemetryConfig             `envconfig:"-"` // This is not loaded from the environment.
	Logging                         LoggingConfig                   `envconfig:"-"` // This is not loaded from the environment.
	Nats                            NatsConfig                      `envconfig:"-"` // This is not loaded from the environment.
	ChatNotificationsConsumerConfig ChatNotificationsConsumerConfig `envconfig:"-"` // This is not loaded from the environment.

	GrpcServerEnabled bool `envconfig:"GRPC_SERVER_ENABLED" default:"true"`
	GrpcServerPort    int  `envconfig:"GRPC_SERVER_PORT" default:"50500"`
}

func (c *Config) Override(cfg Config) {
	if cfg.ServiceName != "" {
		c.ServiceName = cfg.ServiceName
	}
	if cfg.ServiceEnvironment != "" {
		c.ServiceEnvironment = cfg.ServiceEnvironment
	}

	c.Persistence = cfg.Persistence
	c.OpenTelemetry = cfg.OpenTelemetry
	c.Logging = cfg.Logging
	c.Nats = cfg.Nats
	c.ChatNotificationsConsumerConfig = cfg.ChatNotificationsConsumerConfig
	c.GrpcServerEnabled = cfg.GrpcServerEnabled
	c.GrpcServerPort = cfg.GrpcServerPort
}

// MustLoadFromEnv loads the configuration from the environment variables.
// It uses the "RANDOMTALK_MATCHMAKING" prefix for the environment variables.
//
// Global variables are filled and exported to be used in the application.
func MustLoadFromEnv() Config {
	var cfg Config
	envconfig.MustProcess(envPrefix, &cfg)

	cfg.Persistence = mustLoadPersistenceConfig()
	cfg.OpenTelemetry = mustLoadOpenTelemetryConfig()
	cfg.Logging = mustLoadLoggingConfig()
	cfg.Nats = mustLoadNATSConfig()
	cfg.ChatNotificationsConsumerConfig = mustLoadChatNotificationsConsumerConfig()
	return cfg
}
