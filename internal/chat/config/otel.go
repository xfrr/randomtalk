package chatconfig

// Observability holds the configuration for the OpenTelemetry SDK.
type Observability struct {
	// OTELCollectorEndpoint is the endpoint of the OpenTelemetry collector.
	OTELCollectorEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT" default:"http://localhost:4317"`
}
