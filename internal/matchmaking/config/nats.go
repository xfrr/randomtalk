package matchmakingconfig

import "github.com/kelseyhightower/envconfig"

// NatsConfig holds the configuration for the NATS messaging system.
type NatsConfig struct {
	URI string `envconfig:"NATS_URI" default:"nats://localhost:4222"`
}

func mustLoadNATSConfig() NatsConfig {
	var cfg NatsConfig
	envconfig.MustProcess(envPrefix, &cfg)
	return cfg
}
