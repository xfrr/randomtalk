package chatconfig

// NatsConfig holds the configuration for the NATS messaging system.
type NatsConfig struct {
	URI string `envconfig:"NATS_URI" default:"nats://localhost:4222"`
}
