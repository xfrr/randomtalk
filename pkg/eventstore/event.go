package eventstore

import cloudevents "github.com/cloudevents/sdk-go/v2"

const (
	DefaultEventSpecVersion = "1.0"
)

// Event represents a Event according to the CloudEvents Specification.
type Event = cloudevents.Event

type DataContentType string

const (
	ContentTypeApplicationJSON DataContentType = cloudevents.ApplicationJSON
)

func NewEvent() Event {
	evt := cloudevents.NewEvent()
	evt.SetSpecVersion(DefaultEventSpecVersion)
	return evt
}
