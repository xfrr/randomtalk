package eventstore

import (
	"context"
	"time"
)

type SortBy struct {
	Field string
	Order int
}

type FetchOptions struct {
	// Subject filters the events by subject.
	Subject string

	// SortBy sorts the events by time.
	SortBy *SortBy

	// MaxWaitTime is the maximum time to wait for events.
	MaxWaitTime time.Duration
}

type FetchOption func(*FetchOptions)

// FetchSubject sets the subject to filter the events.
func FetchSubject(subject string) FetchOption {
	return func(o *FetchOptions) {
		o.Subject = subject
	}
}

// FetchSortBy sets the sort field and order.
func FetchSortBy(field string, order int) FetchOption {
	return func(o *FetchOptions) {
		o.SortBy = &SortBy{Field: field, Order: order}
	}
}

// FetchMaxWait sets the maximum time to wait for events.
func FetchMaxWait(t time.Duration) FetchOption {
	return func(o *FetchOptions) {
		o.MaxWaitTime = t
	}
}

// AppendResult represents the result of an append operation.
type AppendResult struct {
	// StreamName is the name of the stream.
	StreamName string

	// LastSequence is the sequence number of the last event.
	LastSequence uint64

	// LastEventID is the ID of the last event.
	LastEventID string

	// NumEvents is the number of events appended.
	NumEvents int

	// Error is the error of the operation.
	Error error
}

// Stream exposes the methods to persist and retrieve events.
type Stream interface {
	// Name is the unique name of the Stream.
	Name() string

	// Append adds a list of events to a stream.
	Append(ctx context.Context, events []Event) (AppendResult, error)

	// Pull retrieves a list of events from a stream.
	Pull(ctx context.Context, batchSize int, options ...FetchOption) ([]Event, error)

	// Fetch retrieves a list of events from a stream.
	Fetch(ctx context.Context, batchSize int, options ...FetchOption) (<-chan []Event, error)

	// FetchLast retrieves the last event from a stream.
	FetchLast(ctx context.Context, options ...FetchOption) (*Event, error)
}
