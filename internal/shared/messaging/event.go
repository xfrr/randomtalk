package messaging

import (
	"net/http"
	"sync"

	"github.com/cloudevents/sdk-go/v2/event"
)

// Event is a generic event interface.
type Event struct {
	event.Event

	ackCh   chan struct{}
	ackOnce sync.Once

	nackCh   chan struct{}
	nackOnce sync.Once

	rejectCh   chan struct{}
	rejectOnce sync.Once

	header http.Header
}

// NewEvent creates a new event with the given event and ack/nack functions.
func NewEvent() *Event {
	return &Event{
		Event:    event.New(),
		ackCh:    make(chan struct{}, 1),
		nackCh:   make(chan struct{}, 1),
		rejectCh: make(chan struct{}, 1),
		header:   make(map[string][]string),
	}
}

// Ack acknowledges the event.
func (e *Event) Ack() {
	if e == nil {
		return
	}

	e.ackOnce.Do(func() {
		close(e.ackCh)
	})
}

// Nack negatively acknowledges the event.
func (e *Event) Nack() {
	if e == nil {
		return
	}

	e.nackOnce.Do(func() {
		close(e.nackCh)
	})
}

// Reject rejects the event.
func (e *Event) Reject() {
	if e == nil {
		return
	}

	e.rejectOnce.Do(func() {
		close(e.rejectCh)
	})
}

// Header returns the event header.
func (e *Event) Header() http.Header {
	return e.header
}

// SetHeader sets the event header.
func (e *Event) SetHeader(header map[string][]string) {
	e.header = header
}

// WaitAck waits for the event to be acknowledged.
func (e *Event) WaitAck() <-chan struct{} {
	return e.ackCh
}

// WaitNack waits for the event to be negatively acknowledged.
func (e *Event) WaitNack() <-chan struct{} {
	return e.nackCh
}

// WaitReject waits for the event to be rejected.
func (e *Event) WaitReject() <-chan struct{} {
	return e.rejectCh
}
