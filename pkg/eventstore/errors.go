package eventstore

import "errors"

var (
	// ErrEventNotFound is returned when an event is not found searching for it.
	ErrEventNotFound = errors.New("event not found")

	// ErrNoEventsFound is returned when no events are found.
	ErrNoEventsFound = errors.New("no events found")

	// ErrSequenceMismatch is returned when the sequence of the event is not the expected.
	ErrSequenceMismatch = errors.New("event sequence mismatch")
)
