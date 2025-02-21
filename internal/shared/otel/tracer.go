package xotel

import (
	"context"
	"fmt"
	"time"

	"github.com/xfrr/randomtalk/internal/shared/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0" // Replace with latest semconv version
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type initTracerOptions struct {
	endpointURL        string
	serviceName        string
	serviceVersion     string
	serviceNamespace   string
	serviceEnvironment env.Environment
}

type InitTracerOption func(*initTracerOptions)

// WithServiceName sets the name of the service to be reported in traces.
func WithServiceName(name string) InitTracerOption {
	return func(o *initTracerOptions) {
		o.serviceName = name
	}
}

// WithServiceVersion sets the version of the service to be reported in traces.
func WithServiceVersion(version string) InitTracerOption {
	return func(o *initTracerOptions) {
		o.serviceVersion = version
	}
}

// WithServiceEnvironment sets the environment of the service to be reported in traces.
func WithServiceEnvironment(env env.Environment) InitTracerOption {
	return func(o *initTracerOptions) {
		o.serviceEnvironment = env
	}
}

// WithServiceNamespace sets the namespace of the service to be reported in traces.
func WithServiceNamespace(namespace string) InitTracerOption {
	return func(o *initTracerOptions) {
		o.serviceNamespace = namespace
	}
}

// WithEndpointURL sets the endpoint URL of the OTLP collector.
func WithEndpointURL(url string) InitTracerOption {
	return func(o *initTracerOptions) {
		o.endpointURL = url
	}
}

// initTracer sets up and returns an OpenTelemetry TracerProvider configured with a gRPC exporter.
// It uses environment variables for dynamic configuration and service discovery, and employs
// batch processing for high performance. This function also ensures that resources (e.g. service name)
// are properly identified and annotated.
func InitTracerProvider(ctx context.Context, opts ...InitTracerOption) (*sdktrace.TracerProvider, error) {
	options := initTracerOptions{
		endpointURL:        "localhost:4317",
		serviceName:        "randomtalk.unknown",
		serviceVersion:     "0.1.0",
		serviceEnvironment: "development",
		serviceNamespace:   "randomtalk",
	}

	for _, opt := range opts {
		opt(&options)
	}

	client, err := grpc.NewClient(
		options.endpointURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	// Build the OTLP trace exporter based on the gRPC client.
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(options.serviceName),
			semconv.ServiceVersionKey.String(options.serviceVersion),
			semconv.DeploymentEnvironmentKey.String(options.serviceEnvironment.String()),
			semconv.ServiceNamespaceKey.String(options.serviceNamespace),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Optionally configure a Sampler for controlling tracing volume.
	// e.g., sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)) for 10% sampling.
	//
	// Use BatchSpanProcessor for efficient, asynchronous transport of spans.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			// Fine-tune send parameters for production usage.
			sdktrace.WithMaxQueueSize(2048),
			sdktrace.WithMaxExportBatchSize(512),
			sdktrace.WithBatchTimeout(5*time.Second),
		),
		sdktrace.WithResource(res),
		// Uncomment and configure a sampler to control the rate of tracing.
		// sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.05)),
	)

	// propagate traces using context
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Register our TracerProvider as global so instrumented libraries pick it up automatically.
	otel.SetTracerProvider(tp)

	return tp, nil
}
