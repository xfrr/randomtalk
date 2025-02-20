package matchmakingconfig

import "github.com/kelseyhightower/envconfig"

// OpenTelemetryConfig holds the configuration for the OpenTelemetry SDK.
type OpenTelemetryConfig struct {
	// CollectorEndpoint is the endpoint of the OpenTelemetry collector.
	CollectorEndpoint string `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:"http://localhost:4317"`
}

func mustLoadOpenTelemetryConfig() OpenTelemetryConfig {
	var cfg OpenTelemetryConfig
	envconfig.MustProcess(envPrefix, &cfg)
	return cfg
}
