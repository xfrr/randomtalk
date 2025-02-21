package matchmakingconfig

import "github.com/caarlos0/env/v11"

const envPrefix = "RANDOMTALK_MATCHMAKING_"

type Config struct {
	ServiceName        string `env:"SERVICE_NAME" default:"randomtalk-matchmaking"`
	ServiceEnvironment string `env:"SERVICE_ENVIRONMENT" default:"development"`

	Persistence                     `envPrefix:"PERSISTENCE_"`
	Observability                   `envPrefix:"OBSERVABILITY_"`
	LoggingConfig                   `envPrefix:"LOGGING_"`
	NatsConfig                      `envPrefix:"NATS_"`
	ChatNotificationsConsumerConfig `envPrefix:"CHAT_NOTIFICATIONS_CONSUMER_"`

	GrpcServerEnabled bool `env:"GRPC_SERVER_ENABLED" default:"true"`
	GrpcServerPort    int  `env:"GRPC_SERVER_PORT" default:"50000"`
}

func MustLoadFromEnv() Config {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{
		Prefix:              envPrefix,
		TagName:             "env",
		RequiredIfNoDef:     true,
		DefaultValueTagName: "default",
	})
	if err != nil {
		panic(err)
	}
	return cfg
}
