package matchmakingconfig

import "github.com/kelseyhightower/envconfig"

// LoggingConfig holds the configuration for the logging system.
type LoggingConfig struct {
	Level string `envconfig:"LOG_LEVEL" default:"debug"`
}

func mustLoadLoggingConfig() LoggingConfig {
	var cfg LoggingConfig
	envconfig.MustProcess(envPrefix, &cfg)
	return cfg
}
