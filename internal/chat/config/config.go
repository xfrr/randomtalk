package chatconfig

import (
	"github.com/caarlos0/env/v11"
)

const envPrefix = "RANDOMTALK_CHAT_"

// Config holds the configuration for the randomtalk chat service.
type Config struct {
	ServiceName        string `env:"SERVICE_NAME" default:"randomtalk-chat"`
	ServiceEnvironment string `env:"SERVICE_ENVIRONMENT" default:"development"`

	MatchNotificationsConsumerConfig `envPrefix:"NATS_MATCH_NOTIFICATIONS_CONSUMER_"`
	ChatSessionStreamConfig          `envPrefix:"CHAT_SESSION_STREAM_"`
	NotificationStreamConfig         `envPrefix:"NATS_NOTIFICATION_STREAM_"`
	HubWebsocketServer               `envPrefix:"HUB_WEBSOCKET_SERVER_"`
	LoggingConfig                    `envPrefix:"LOGGING_"`
	NatsConfig                       `envPrefix:"NATS_"`
	Observability                    `envPrefix:"OBSERVABILITY_"`
}

// MustLoadFromEnv loads the configuration from the environment variables.
// It uses the "RANDOMTALK_CHAT" prefix for the environment variables.
//
// Global variables are filled and exported to be used in the application.
func MustLoadFromEnv() Config {
	var cfg Config
	err := env.ParseWithOptions(&cfg, env.Options{
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
