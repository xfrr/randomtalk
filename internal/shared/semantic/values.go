package semantic

// Message types are used to categorize different kinds of messages
// in a messaging system, such as commands, events, notifications, errors,
// system messages, and unknown types. These constants can be used to
// identify the type of message being processed or transmitted, allowing
// for appropriate handling and processing logic based on the message type.
const (
	MessageTypeCommand      = "command"
	MessageTypeEvent        = "event"
	MessageTypeNotification = "notification"
	MessageTypeError        = "error"
	MessageTypeSystem       = "system"
	MessageTypeUnknown      = "unknown"
)
