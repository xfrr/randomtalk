package matchmakingconfig

// NatsConfig holds the configuration for the NATS messaging system.
type NatsConfig struct {
	URI string `env:"URI" default:"nats://localhost:4222"`
}
