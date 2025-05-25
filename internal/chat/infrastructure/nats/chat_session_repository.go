package chatnats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/aggregate/event"
	chatdom "github.com/xfrr/randomtalk/internal/chat/domain"
	chatdomaineventsv1 "github.com/xfrr/randomtalk/internal/chat/domain/events/v1"
	"github.com/xfrr/randomtalk/internal/shared/eventstore"
	xnats "github.com/xfrr/randomtalk/internal/shared/nats"
)

const (
	chatSessionsStreamSuffix = "sessions"
)

// ensure MatchRepository implements chatdom.MatchRepository
var _ chatdom.ChatSessionRepository = (*ChatSessionRepository)(nil)

// MatchRepository implements chatdom.MatchRepository using NATS JetStream.
type ChatSessionRepository struct {
	sourceName string
	stream     eventstore.Stream
}

func NewChatSessionRepository(ctx context.Context, js jetstream.JetStream, streamConfig xnats.StreamConfig) (*ChatSessionRepository, error) {
	stream, err := xnats.CreateStream(ctx, js, streamConfig)
	if err != nil {
		return nil, err
	}

	return &ChatSessionRepository{
		sourceName: buildStreamSourceName(chatdom.EventSourceName, chatSessionsStreamSuffix),
		stream:     stream,
	}, nil
}

// Save appends new events for a Match to the event store with optimistic concurrency checks.
func (r ChatSessionRepository) Save(ctx context.Context, match *chatdom.ChatSession) error {
	events, err := r.toStoreEvents(match.AggregateEvents())
	if err != nil {
		return fmt.Errorf("convert to cloud events: %w", err)
	}

	_, appendErr := r.stream.Append(ctx, events)
	if appendErr != nil {
		if errors.Is(appendErr, eventstore.ErrSequenceMismatch) {
			return chatdom.ErrChatSessionAlreadyExists
		}
		return fmt.Errorf("append chat session events: %w", appendErr)
	}

	return nil
}

func (r ChatSessionRepository) FindByID(ctx context.Context, id string) (*chatdom.ChatSession, error) {
	batchSize := 5
	filter := eventstore.FetchSubject(createEventFilterKey(r.sourceName, id, ">"))

	cloudEvents, err := r.stream.Pull(ctx, batchSize, filter)
	if err != nil {
		if errors.Is(err, eventstore.ErrNoEventsFound) {
			return nil, chatdom.ErrChatSessionNotFound
		}
		return nil, fmt.Errorf("pull from stream: %w", err)
	}

	aggEvents, err := eventsFromCloudEvents(cloudEvents)
	if err != nil {
		return nil, fmt.Errorf("convert from cloud events: %w", err)
	}

	sess, restoreErr := chatdom.NewChatSessionFromEvents(chatdom.ID(id), aggEvents)
	if restoreErr != nil {
		return nil, fmt.Errorf("restore chat session from events: %w", restoreErr)
	}
	return sess, nil
}

func (r ChatSessionRepository) Exists(ctx context.Context, id string) (bool, error) {
	filter := eventstore.FetchSubject(createEventFilterKey(r.sourceName, chatdom.AggregateName, id, ">"))
	cloudEvents, err := r.stream.Pull(ctx, 1, filter)
	if err != nil {
		if errors.Is(err, eventstore.ErrNoEventsFound) {
			return false, nil
		}
		return false, fmt.Errorf("pull from stream (Exists): %w", err)
	}

	return len(cloudEvents) > 0, nil
}

func (r ChatSessionRepository) toStoreEvents(events []aggregate.Event) ([]eventstore.Event, error) {
	cloudEvents := make([]eventstore.Event, 0, len(events))
	for _, evt := range events {
		aggregateID, ok := evt.Aggregate().ID.(string)
		if !ok {
			return nil, errors.New("aggregate ID must be a string")
		}
		eventID, ok := evt.ID().(string)
		if !ok {
			return nil, errors.New("event ID must be a string")
		}

		ce := eventstore.NewEvent()
		ce.SetID(eventID)
		ce.SetType(evt.Name())
		ce.SetSource(chatdom.EventSourceName)
		ce.SetSubject(strings.Join([]string{chatSessionsStreamSuffix, aggregateID}, "."))
		ce.SetTime(evt.OccurredAt())
		ce.SetDataSchema("schemas.randomtalk.com/chat/events/" + evt.Name() + "/1.0")

		if err := ce.Context.SetExtension(xnats.SubjectVersionHeaderKey, strconv.Itoa(evt.Aggregate().Version)); err != nil {
			return nil, fmt.Errorf("set extension: %w", err)
		}
		if dataErr := ce.SetData(string(eventstore.ContentTypeApplicationJSON), evt.Payload()); dataErr != nil {
			return nil, fmt.Errorf("set event data: %w", dataErr)
		}
		cloudEvents = append(cloudEvents, ce)
	}
	return cloudEvents, nil
}

func buildStreamSourceName(sourceName, streamSuffix string) string {
	return sourceName + "." + streamSuffix
}

func eventsFromCloudEvents(cloudEvents []eventstore.Event) ([]aggregate.Event, error) {
	aggEvents := make([]aggregate.Event, len(cloudEvents))
	for i, ce := range cloudEvents {
		evt, err := eventFromCloudEvent(ce)
		if err != nil {
			return nil, fmt.Errorf("convert single cloud event: %w", err)
		}
		aggEvents[i] = evt
	}
	return aggEvents, nil
}

func eventFromCloudEvent(ce eventstore.Event) (aggregate.Event, error) {
	aggVersion, err := types.ToInteger(xnats.SubjectVersionFromMap(ce.Extensions()))
	if err != nil {
		return nil, fmt.Errorf("invalid event aggregate version: %w", err)
	}

	switch ce.Type() {
	case chatdomaineventsv1.ChatSessionCreated{}.EventName():
		payload := &chatdomaineventsv1.ChatSessionCreated{}
		if unmarshalErr := json.Unmarshal(ce.DataEncoded, payload); unmarshalErr != nil {
			return nil, fmt.Errorf("json unmarshal: %w", unmarshalErr)
		}

		subjectSplit := strings.Split(ce.Subject(), ".")
		if len(subjectSplit) < 2 {
			return nil, errors.New("invalid subject format")
		}

		subjectID := subjectSplit[1]
		if payload.UserID != subjectID {
			return nil, errors.New("subject ID and event payload ID mismatch")
		}

		ev, evErr := event.New(
			ce.ID(),
			ce.Type(),
			payload,
			event.WithOccurredAt(ce.Time()),
			event.WithAggregate(
				subjectID,
				chatdom.AggregateName,
				int(aggVersion),
			),
		)
		if evErr != nil {
			return nil, fmt.Errorf("create domain event: %w", evErr)
		}
		return ev.Any(), nil
	default:
		return nil, fmt.Errorf("unexpected event type: %s", ce.Type())
	}
}

func createEventFilterKey(sourceName string, parts ...string) string {
	return strings.Join(append([]string{sourceName}, parts...), ".")
}
