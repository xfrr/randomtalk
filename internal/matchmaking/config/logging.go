package matchmakingconfig

// LoggingConfig holds the configuration for the logging system.
type LoggingConfig struct {
	Level string `env:"LEVEL" default:"debug"`
}
