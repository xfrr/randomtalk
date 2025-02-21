package chatconfig

import "github.com/kelseyhightower/envconfig"

// Config holds the configuration for the randomtalk chat service.
type Config struct {
	ServiceName        string `envconfig:"SERVICE_NAME" default:"randomtalk-chat"`
	ServiceEnvironment string `envconfig:"SERVICE_ENVIRONMENT" default:"development"`

	// NotificationStream is the configuration for the JetStream stream used for notifications.
	NotificationStream NotificationStreamConfig
	// WebsocketServer is the configuration for the websocket server.
	WebsocketServer HubWebsocketServer
	// Logging is the configuration for the logging system.
	Logging LoggingConfig
	// Nats is the configuration for the NATS server.
	Nats NatsConfig
	// OpenTelemetry is the configuration for the OpenTelemetry system.
	OpenTelemetry OpenTelemetryConfig
}

// MustLoadFromEnv loads the configuration from the environment variables.
// It uses the "RANDOMTALK_CHAT" prefix for the environment variables.
//
// Global variables are filled and exported to be used in the application.
func MustLoadFromEnv() Config {
	var cfg Config
	envconfig.MustProcess(envPrefix, &cfg)
	return cfg
}
