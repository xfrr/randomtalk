package matchnats

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
	"github.com/xfrr/randomtalk/pkg/eventstore"

	matchdom "github.com/xfrr/randomtalk/internal/matchmaking/domain"
)

const (
	streamEventAggregateVersionHeaderKey = "subjectversion"
	matchesStreamSuffix                  = "matches"
)

// ensure MatchRepository implements matchdom.MatchRepository
var _ matchdom.MatchRepository = (*MatchRepository)(nil)

// MatchRepository implements matchdom.MatchRepository using NATS JetStream.
type MatchRepository struct {
	sourceName string
	stream     eventstore.Stream
}

func NewMatchStreamRepository(ctx context.Context, js jetstream.JetStream, streamConfig StreamConfig) (*MatchRepository, error) {
	// create nats match events stream
	stream, err := CreateStream(ctx, js, streamConfig)
	if err != nil {
		return nil, err
	}

	return &MatchRepository{
		sourceName: buildStreamSourceName(matchdom.EventSourceName, matchesStreamSuffix),
		stream:     stream,
	}, nil
}

// Save appends new events for a Match to the event store with optimistic concurrency checks.
func (r *MatchRepository) Save(ctx context.Context, match *matchdom.Match) error {
	events, err := r.toCloudEvents(match.AggregateEvents())
	if err != nil {
		return fmt.Errorf("convert to cloud events: %w", err)
	}

	res, appendErr := r.stream.Append(ctx, events)
	if appendErr != nil {
		if errors.Is(appendErr, eventstore.ErrSequenceMismatch) {
			if res.LastSequence == 0 {
				return matchdom.ErrMatchAlreadyExists
			}
			return matchdom.ErrMatchAlreadyExists
		}
		return fmt.Errorf("append match events: %w", appendErr)
	}
	return nil
}

func (r *MatchRepository) FindByID(ctx context.Context, id string) (*matchdom.Match, error) {
	batchSize := 5
	filter := eventstore.FetchSubject(createEventFilterKey(r.sourceName, id, ">"))

	cloudEvents, err := r.stream.Pull(ctx, batchSize, filter)
	if err != nil {
		if errors.Is(err, eventstore.ErrNoEventsFound) {
			return nil, matchdom.ErrMatchNotFound
		}
		return nil, fmt.Errorf("pull from stream: %w", err)
	}

	aggEvents, err := eventsFromCloudEvents(cloudEvents)
	if err != nil {
		return nil, fmt.Errorf("convert from cloud events: %w", err)
	}

	sess, restoreErr := matchdom.NewMatchFromEvents(matchdom.MatchID(id), aggEvents...)
	if restoreErr != nil {
		return nil, fmt.Errorf("restore match from events: %w", restoreErr)
	}
	return sess, nil
}

func (r *MatchRepository) Exists(ctx context.Context, id string) (bool, error) {
	filter := eventstore.FetchSubject(createEventFilterKey(r.sourceName, id, ">"))
	cloudEvents, err := r.stream.Pull(ctx, 1, filter)
	if err != nil {
		if errors.Is(err, eventstore.ErrNoEventsFound) {
			return false, nil
		}
		return false, fmt.Errorf("pull from stream (Exists): %w", err)
	}
	return len(cloudEvents) > 0, nil
}

func (r *MatchRepository) FindLastByUserID(ctx context.Context, userID string) (*matchdom.Match, error) {
	// 1. retrieve all match events
	batchSize := 10
	filter := createEventFilterKey(r.sourceName, matchdom.MatchAggregateName, "*.created")

	// 2. iterate over events and filter by user ID
	cloudEvents, err := r.stream.Fetch(ctx, batchSize, eventstore.FetchSubject(filter))
	if err != nil {
		if errors.Is(err, eventstore.ErrNoEventsFound) {
			return nil, matchdom.ErrMatchNotFound
		}

		return nil, fmt.Errorf("fetch from stream: %w", err)
	}

	for cevents := range cloudEvents {
		aggEvents, err := eventsFromCloudEvents(cevents)
		if err != nil {
			return nil, fmt.Errorf("convert from cloud events: %w", err)
		}

		if len(aggEvents) == 0 {
			return nil, matchdom.ErrMatchNotFound
		}

		aggID, ok := aggEvents[0].Aggregate().ID.(string)
		if !ok {
			return nil, errors.New("aggregate ID must be a string")
		}

		match, restoreErr := matchdom.NewMatchFromEvents(matchdom.MatchID(aggID), aggEvents...)
		if restoreErr != nil {
			return nil, fmt.Errorf("restore match from events: %w", restoreErr)
		}

		if match.Requester().ID() == userID || match.Candidate().ID() == userID {
			return match, nil
		}
	}

	return nil, matchdom.ErrMatchNotFound
}

func (r *MatchRepository) toCloudEvents(events []aggregate.Event) ([]eventstore.Event, error) {
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
		ce.SetSource(matchdom.EventSourceName)
		ce.SetSubject(strings.Join([]string{"matches", aggregateID}, "."))
		ce.SetTime(evt.OccurredAt())
		ce.SetDataSchema("schemas.randomtalk.com/matchmaking/match/events/" + evt.Name() + "/1.0")

		if err := ce.Context.SetExtension(streamEventAggregateVersionHeaderKey, strconv.Itoa(evt.Aggregate().Version)); err != nil {
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
	aggVersion, err := types.ToInteger(ce.Extensions()[streamEventAggregateVersionHeaderKey])
	if err != nil {
		return nil, fmt.Errorf("invalid event aggregate version: %w", err)
	}

	switch ce.Type() {
	case matchdom.MatchCreatedEvent{}.EventName():
		payload := &matchdom.MatchCreatedEvent{}
		if unmarshalErr := json.Unmarshal(ce.DataEncoded, payload); unmarshalErr != nil {
			return nil, fmt.Errorf("json unmarshal: %w", unmarshalErr)
		}

		subjectSplit := strings.Split(ce.Subject(), ".")
		if len(subjectSplit) < 2 {
			return nil, errors.New("invalid subject format")
		}

		subjectID := subjectSplit[1]
		if payload.MatchID != subjectID {
			return nil, errors.New("subject ID and event payload ID mismatch")
		}

		ev, evErr := event.New(
			ce.ID(),
			ce.Type(),
			payload,
			event.WithOccurredAt(ce.Time()),
			event.WithAggregate(
				subjectID,
				matchdom.MatchAggregateName,
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

func createEventFilterKey(sourceName, id string, eventType ...string) string {
	base := sourceName + "." + id
	if len(eventType) == 0 {
		return base
	}
	return base + "." + eventType[0]
}
