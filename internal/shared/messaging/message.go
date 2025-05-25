package messaging

// Message represents a message that can be sent or received.
type Message interface {
	// ID returns the message ID.
	ID() string

	// Type returns the message type.
	Type() string

	// Payload returns the message payload.
	Payload() []byte
}
