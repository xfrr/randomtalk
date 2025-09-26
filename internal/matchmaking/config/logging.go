package matchmakingconfig

// LoggingConfig holds the configuration for the logging system.
type LoggingConfig struct {
	LogLevel string `env:"LEVEL" default:"debug"`
}
