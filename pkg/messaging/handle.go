package messaging

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// EventSubscriber is an interface for subscribing to events.
type EventSubscriber interface {
	// Subscribe returns a channel to receive events.
	Subscribe(context.Context) (<-chan *Event, error)
}

// EventHandlerFunc is a function that handles a messaging event.
type EventHandlerFunc func(context.Context, *Event) error

// HandleEvents subscribes to events and handles them sequentially in the same goroutine.
func HandleEvents(ctx context.Context, log *zerolog.Logger, subs EventSubscriber, handlerFn EventHandlerFunc) error {
	events, err := subs.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	for evt := range events {
		if evt == nil {
			continue
		}

		sbctx := extractSpanContext(context.Background(), evt)
		span := trace.SpanFromContext(sbctx)
		span.SetName(evt.Type())
		span.SetAttributes(
			attribute.String("event_id", evt.ID()),
			attribute.String("event_type", evt.Type()),
		)

		// call the handler passing the span context and the event
		if handleErr := handlerFn(sbctx, evt); handleErr != nil {
			span.RecordError(handleErr)
			span.SetStatus(codes.Error, handleErr.Error())
			log.Error().
				Err(handleErr).
				Str("event_id", evt.ID()).
				Msgf("failed to handle %s event", evt.Type())
		}

		span.End()
	}

	return nil
}

// HandleEventsInWorkerPool subscribes to events and processes them using numWorkers goroutines.
// Each worker reads events from the subscription channel, applies the same tracing logic,
// and calls the provided handler function in parallel.
func HandleEventsInWorkerPool(
	ctx context.Context,
	log zerolog.Logger,
	subs EventSubscriber,
	handlerFn EventHandlerFunc,
	numWorkers int,
) error {
	if numWorkers < 1 {
		numWorkers = 1
	}

	events, err := subs.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					// Context canceled; exit worker
					return
				case evt, ok := <-events:
					if !ok {
						// Channel closed; exit worker
						return
					}
					if evt == nil {
						continue
					}

					// Extract tracing context
					sbctx := extractSpanContext(context.Background(), evt)
					span := trace.SpanFromContext(sbctx)
					span.SetName(evt.Type())
					span.SetAttributes(
						attribute.String("event_id", evt.ID()),
						attribute.String("event_type", evt.Type()),
						attribute.Int("worker_id", workerID),
					)

					// call the handler passing the span context and the event
					if handleErr := handlerFn(sbctx, evt); handleErr != nil {
						span.RecordError(handleErr)
						span.SetStatus(codes.Error, handleErr.Error())
						log.Error().
							Err(handleErr).
							Str("event_id", evt.ID()).
							Int("worker_id", workerID).
							Msgf("failed to handle %s event", evt.Type())
					}
					span.End()
				}
			}
		}(i)
	}

	// Wait until all workers are done (either ctx canceled or channel closed).
	wg.Wait()
	return ctx.Err()
}

// extractSpanContext propagates the tracing headers from event metadata into the current context.
func extractSpanContext(ctx context.Context, msg *Event) context.Context {
	propagator := propagation.TraceContext{}
	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(msg.Header()))
	return ctx
}
