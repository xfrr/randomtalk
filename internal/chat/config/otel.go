package chatconfig

// OpenTelemetryConfig holds the configuration for the OpenTelemetry SDK.
type OpenTelemetryConfig struct {
	CollectorEndpoint string `envconfig:"OTEL_COLLECTOR_ENDPOINT" default:"http://localhost:4317"`
}
