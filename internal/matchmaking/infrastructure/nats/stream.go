package matchnats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/randomtalk/pkg/eventstore"
)

var _ eventstore.Stream = (*Stream)(nil)

// Stream is an event store backed by NATS JetStream.
type Stream struct {
	streamConfig jetstream.StreamConfig
	js           jetstream.JetStream
	stream       jetstream.Stream
}

// CreateStream creates or updates a JetStream stream and returns a Stream implementation.
func CreateStream(ctx context.Context, js jetstream.JetStream, config StreamConfig) (*Stream, error) {
	stream, err := js.CreateOrUpdateStream(ctx, config.streamConfig)
	if err != nil {
		return nil, err
	}

	return &Stream{
		streamConfig: config.streamConfig,
		js:           js,
		stream:       stream,
	}, nil
}

// Name returns the JetStream stream name.
func (s *Stream) Name() string {
	return s.streamConfig.Name
}

// Append publishes the provided events to the stream.
func (s *Stream) Append(ctx context.Context, events []event.Event) (eventstore.AppendResult, error) {
	res := eventstore.AppendResult{
		StreamName: s.streamConfig.Name,
	}
	if len(events) == 0 {
		return res, errors.New("no events to append")
	}

	// Determine last sequence for concurrency checks (if needed).
	firstEvent := events[0]
	natsSubject := s.makeSubjectFromEvent(firstEvent)

	var err error
	res.LastEventID, res.LastSequence, err = s.getLastSequence(ctx, natsSubject)
	if err != nil {
		return res, err
	}

	for _, e := range events {
		encoded, encodingError := e.MarshalJSON()
		if encodingError != nil {
			return res, encodingError
		}

		// Basic publish options, including a message ID for deduplication & retries.
		publishOpts := []jetstream.PublishOpt{
			jetstream.WithMsgID(e.ID()),
			jetstream.WithRetryAttempts(3),
			jetstream.WithRetryWait(200 * time.Millisecond),
			jetstream.WithExpectStream(s.streamConfig.Name),
			// Concurrency checks: using LastMsgID for this example
			// jetstream.WithExpectLastMsgID(res.LastEventID),
			// or if you'd prefer strict subject-level concurrency:
			// jetstream.WithExpectLastSequencePerSubject(res.LastSequence),
		}

		// Retrieve aggregate version from CloudEvents extensions, if present.
		aggregateID, extErr := types.ToString(e.Extensions()[streamEventAggregateVersionHeaderKey])
		if extErr != nil {
			// If no aggregate version is found, you could skip or treat this as a zero version.
			// For now, we'll return an error to preserve existing behavior.
			return res, extErr
		}

		publishSubject := s.makeSubjectFromEvent(e)
		natsMsg := &nats.Msg{
			Subject: publishSubject,
			Data:    encoded,
			Header: nats.Header{
				"Content-Type":                       []string{"application/cloudevents+json"},
				streamEventAggregateVersionHeaderKey: []string{aggregateID},
			},
		}

		puback, pubErr := s.js.PublishMsg(ctx, natsMsg, publishOpts...)
		if pubErr != nil {
			switch {
			case errors.Is(pubErr, jetstream.ErrKeyExists):
				// KeyExists indicates a message with this ID already exists (sequence mismatch).
				return res, eventstore.ErrSequenceMismatch
			default:
				return res, pubErr
			}
		}

		// Update last sequence and event ID after a successful publish.
		res.LastSequence = puback.Sequence
		res.LastEventID = e.ID()
	}

	return res, nil
}

// Pull fetches a batch of events from a new ephemeral consumer, returning them as a slice.
func (s *Stream) Pull(ctx context.Context, batchSize int, options ...eventstore.FetchOption) ([]eventstore.Event, error) {
	fetchOptions := &eventstore.FetchOptions{}
	for _, opt := range options {
		opt(fetchOptions)
	}

	consumer, err := s.newOrderedConsumer(ctx, fetchOptions.Subject)
	if err != nil {
		return nil, err
	}

	// If batchSize <= 0, fallback to total messages in the stream.
	bsize := batchSize
	if bsize <= 0 {
		bsize, err = s.fetchBatchSize(ctx)
		if err != nil {
			return nil, err
		}
	}
	if bsize == 0 {
		return nil, nil
	}

	messages, err := consumer.FetchNoWait(bsize)
	if err != nil {
		switch {
		case errors.Is(err, jetstream.ErrNoMessages):
			return nil, eventstore.ErrNoEventsFound
		default:
			return nil, err
		}
	}

	var events []eventstore.Event
	for msg := range messages.Messages() {
		e := eventstore.NewEvent()
		if unmarshalErr := e.UnmarshalJSON(msg.Data()); unmarshalErr != nil {
			return nil, unmarshalErr
		}
		events = append(events, e)
	}

	if len(events) == 0 {
		return nil, eventstore.ErrNoEventsFound
	}
	return events, nil
}

// Fetch streams events continuously in batches, sending them to the returned channel.
func (s *Stream) Fetch(ctx context.Context, batchSize int, options ...eventstore.FetchOption) (<-chan []eventstore.Event, error) {
	fetchOptions := &eventstore.FetchOptions{
		// TODO: change to 30s
		MaxWaitTime: 10 * time.Second,
	}
	for _, opt := range options {
		opt(fetchOptions)
	}

	consumer, err := s.newOrderedConsumer(ctx, fetchOptions.Subject)
	if err != nil {
		return nil, err
	}

	eventsCh := make(chan []eventstore.Event)
	go func() {
		defer close(eventsCh)

		for {
			// Respect caller cancellation.
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Fetch the next batch.
			msgBatch, fetchErr := consumer.Fetch(batchSize, jetstream.FetchMaxWait(fetchOptions.MaxWaitTime))
			if fetchErr != nil {
				if errors.Is(fetchErr, jetstream.ErrNoMessages) {
					// No messages found, continue fetching to keep streaming.
					continue
				}
				// Could log or handle other errors, then exit.
				return
			}

			// Decode messages.
			msgsDecoded := make([]eventstore.Event, 0, batchSize)
			for msg := range msgBatch.Messages() {
				e, decodeErr := decodeEvent(msg.Data())
				if decodeErr != nil {
					// Could log decodeErr, then skip.
					// For simplicity, we'll just return to close the goroutine.
					continue
				}

				msgsDecoded = append(msgsDecoded, *e)
			}

			eventsCh <- msgsDecoded
		}
	}()

	return eventsCh, nil
}

// FetchLast returns the last published event, optionally filtered by subject.
func (s *Stream) FetchLast(ctx context.Context, options ...eventstore.FetchOption) (*eventstore.Event, error) {
	opts := &eventstore.FetchOptions{}
	for _, opt := range options {
		opt(opts)
	}

	lastEvent, err := s.fetchLastMessage(ctx, opts)
	if err != nil {
		switch {
		case errors.Is(err, jetstream.ErrMsgNotFound):
			return nil, eventstore.ErrEventNotFound
		default:
			return nil, err
		}
	}
	return lastEvent, nil
}

// -----------------------------------------------------------------------------
// Private Helpers
// -----------------------------------------------------------------------------

// newOrderedConsumer abstracts ephemeral consumer creation for Pull/Fetch.
func (s *Stream) newOrderedConsumer(ctx context.Context, subject string) (jetstream.Consumer, error) {
	cfg := jetstream.OrderedConsumerConfig{
		DeliverPolicy: jetstream.DeliverAllPolicy,
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
	}
	if subject != "" {
		cfg.FilterSubjects = []string{subject}
	}
	return s.js.OrderedConsumer(ctx, s.streamConfig.Name, cfg)
}

func (s *Stream) getLastSequence(ctx context.Context, subject string) (string, uint64, error) {
	// Called under Append’s lock, so no extra locking needed here.
	lastMsg, err := s.stream.GetLastMsgForSubject(ctx, subject)
	if err != nil {
		switch {
		case errors.Is(err, jetstream.ErrMsgNotFound):
			return "", 0, nil
		default:
			return "", 0, err
		}
	}
	msgID := lastMsg.Header.Get(nats.MsgIdHdr)
	return msgID, lastMsg.Sequence, nil
}

func (s *Stream) fetchLastMessage(ctx context.Context, opts *eventstore.FetchOptions) (*eventstore.Event, error) {
	// Called under FetchLast’s read lock, no extra locking needed here.
	stream, err := s.js.Stream(ctx, s.streamConfig.Name)
	if err != nil {
		return nil, err
	}

	if opts.Subject != "" {
		return s.fetchLastMessageBySubject(ctx, stream, opts.Subject)
	}
	return s.fetchLastMessageInStream(ctx, stream)
}

func (s *Stream) fetchLastMessageBySubject(ctx context.Context, stream jetstream.Stream, subject string) (*eventstore.Event, error) {
	lastMsg, err := stream.GetLastMsgForSubject(ctx, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get last message for subject %s: %w", subject, err)
	}
	return decodeEvent(lastMsg.Data)
}

func (s *Stream) fetchLastMessageInStream(ctx context.Context, stream jetstream.Stream) (*eventstore.Event, error) {
	info := stream.CachedInfo()
	lastSeq := info.State.LastSeq
	if lastSeq == 0 {
		return nil, jetstream.ErrMsgNotFound
	}

	lastMsg, err := stream.GetMsg(ctx, lastSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to get last message in stream: %w with sequence %d", err, lastSeq)
	}
	return decodeEvent(lastMsg.Data)
}

// fetchBatchSize returns the total messages in the stream, safely cast to int.
func (s *Stream) fetchBatchSize(ctx context.Context) (int, error) {
	// Called under Pull’s lock, no extra locking needed here.
	stream, err := s.js.Stream(ctx, s.streamConfig.Name)
	if err != nil {
		return 0, err
	}
	nmsgs := stream.CachedInfo().State.Msgs
	if nmsgs == 0 {
		return 0, nil
	}
	// Carefully convert uint64 to int.
	if nmsgs < 1<<31 {
		return int(nmsgs), nil
	}
	return 0, errors.New("number of messages in stream exceeds int32")
}

// decodeEvent converts raw bytes to an eventstore.Event.
func decodeEvent(data []byte) (*eventstore.Event, error) {
	e := event.New()
	if err := e.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	ev := eventstore.Event(e) // Convert to eventstore.Event if needed
	return &ev, nil
}

// makeSubjectFromEvent constructs the subject from the event’s source, subject, and type.
func (s *Stream) makeSubjectFromEvent(e event.Event) string {
	return e.Source() + "." + e.Subject() + "." + e.Type()
}
