package eventstoretraces

import (
	"context"

	"github.com/xfrr/randomtalk/pkg/eventstore"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type StreamOpentelemetryMiddleware struct {
	engine eventstore.PersistenceEngine
	stream eventstore.Stream
	tracer trace.Tracer
}

func NewStreamOpentelemetryMiddleware(
	engine eventstore.PersistenceEngine,
	stream eventstore.Stream,
	tracer trace.Tracer,
) eventstore.Stream {
	return &StreamOpentelemetryMiddleware{
		engine: engine,
		stream: stream,
		tracer: tracer,
	}
}

func (s *StreamOpentelemetryMiddleware) Name() string {
	return s.stream.Name()
}

func (s *StreamOpentelemetryMiddleware) FetchLast(ctx context.Context, opts ...eventstore.FetchOption) (*eventstore.Event, error) {
	ctx, span := s.tracer.Start(ctx, "eventstore.stream.fetch_last")
	defer span.End()

	options := eventstore.FetchOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	attrs := []attribute.KeyValue{
		semconv.MessagingDestinationKindTopic,
		attribute.String(string(semconv.MessagingSystemKey), string(s.engine)),
	}

	if options.Subject != "" {
		attrs = append(attrs, attribute.String("eventstore.stream.subject", options.Subject))
	}

	span.SetAttributes(attrs...)
	return s.stream.FetchLast(ctx, opts...)
}

func (s *StreamOpentelemetryMiddleware) Pull(
	ctx context.Context,
	batchSize int,
	opts ...eventstore.FetchOption,
) ([]eventstore.Event, error) {
	ctx, span := s.tracer.Start(ctx, "eventstore.stream.fetch")
	defer span.End()

	options := &eventstore.FetchOptions{}
	for _, opt := range opts {
		opt(options)
	}

	attrs := []attribute.KeyValue{
		semconv.MessagingDestinationKindTopic,
		attribute.String(string(semconv.MessagingSystemKey), string(s.engine)),
		attribute.Int("eventstore.stream.fetch.batch_size", batchSize),
	}

	if options.Subject != "" {
		attrs = append(attrs, attribute.String("eventstore.stream.subject", options.Subject))
	}

	span.SetAttributes(attrs...)
	return s.stream.Pull(ctx, batchSize, opts...)
}

func (s *StreamOpentelemetryMiddleware) Fetch(
	ctx context.Context,
	batchSize int,
	opts ...eventstore.FetchOption,
) (<-chan []eventstore.Event, error) {
	ctx, span := s.tracer.Start(ctx, "eventstore.stream.fetch")
	defer span.End()

	attrs := []attribute.KeyValue{
		semconv.MessagingDestinationKindTopic,
		attribute.String(string(semconv.MessagingSystemKey), string(s.engine)),
		attribute.Int("eventstore.stream.fetch.batch_size", batchSize),
	}

	span.SetAttributes(attrs...)
	return s.stream.Fetch(ctx, batchSize, opts...)
}

func (s *StreamOpentelemetryMiddleware) Append(
	ctx context.Context,
	events []eventstore.Event,
) (eventstore.AppendResult, error) {
	ctx, span := s.tracer.Start(ctx, "eventstore.stream.append")
	defer span.End()

	attrs := []attribute.KeyValue{
		semconv.MessagingDestinationKindTopic,
		attribute.String(string(semconv.MessagingSystemKey), string(s.engine)),
		attribute.Int("eventstore.stream.events_count", len(events)),
	}

	if len(events) > 0 {
		attrs = append(attrs, attribute.String("eventstore.stream.source", events[0].Source()))
		attrs = append(attrs, attribute.String("eventstore.stream.subject", events[0].Subject()))
	}

	span.SetAttributes(attrs...)
	return s.stream.Append(ctx, events)
}
