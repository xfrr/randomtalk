//go:build integration
// +build integration

package xnats_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	eventstore "github.com/xfrr/randomtalk/internal/shared/eventstore"
	xnats "github.com/xfrr/randomtalk/internal/shared/nats"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStream(
	t *testing.T,
	nc *nats.Conn,
	streamName string,
	subjects []string,
) (context.Context, jetstream.JetStream, *xnats.Stream) {
	t.Helper()

	ctx := context.Background()

	js, err := jetstream.New(nc)
	require.NoError(t, err, "Failed to create JetStream context")

	sut, err := xnats.CreateStream(ctx, js, xnats.NewStreamConfig(streamName, subjects...))
	require.NoError(t, err)

	// sleep to ensure the stream is created
	time.Sleep(500 * time.Millisecond)

	stream, err := js.Stream(ctx, streamName)
	require.NoError(t, err, "Failed to get stream")
	require.NotNil(t, stream, "Stream should not be nil")

	// Clean up resources after the test.
	t.Cleanup(func() {
		_ = js.DeleteStream(ctx, streamName)
	})

	_ = stream.Purge(ctx)
	return ctx, js, sut
}

func makeEvent(t *testing.T, id string, source string) cloudevents.Event {
	t.Helper()
	event := cloudevents.NewEvent()
	event.SetID(fmt.Sprintf("test-event-%s", id))
	event.SetType("event_created")
	event.SetSource(source)
	event.SetSubject("test.event")
	event.SetDataContentType("application/json")
	err := event.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"event": "test",
	})
	require.NoError(t, err)

	return event
}

func makeEvents(t *testing.T, n int, source string, prefix ...string) []eventstore.Event {
	t.Helper()
	events := make([]eventstore.Event, n)
	for i := 0; i < n; i++ {
		eventID := fmt.Sprintf("%s-%d", strings.Join(prefix, ""), i)
		events[i] = makeEvent(t, eventID, source)
	}
	return events
}

func TestEventStoreStream_Append(t *testing.T) {
	streamName := "TEST_EVENTSTORE_STREAM_APPEND"

	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err, "Failed to connect to NATS")
	defer nc.Close()

	tests := []struct {
		name              string
		events            []eventstore.Event
		expectedNumEvents int
		expectedLastSeq   uint64
		expectErr         bool
		expectedErrMsg    string
		expectedErrType   error
	}{
		{
			name:              "Append with empty events list",
			events:            []eventstore.Event{},
			expectedNumEvents: 0,
			expectedLastSeq:   0,
			expectErr:         false,
		},
		{
			name:              "Append a list of events",
			events:            makeEvents(t, 5, "test.append"),
			expectedNumEvents: 5,
			expectedLastSeq:   5,
			expectErr:         false,
		},
	}

	for _, test := range tests {
		tt := test // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, js, sut := setupTestStream(t, nc, streamName, []string{"test.append.>"})

			result, err := sut.Append(ctx, tt.events)
			if tt.expectErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.IsType(t, tt.expectedErrType, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedNumEvents, result.NumEvents)
			require.Equal(t, int(tt.expectedLastSeq), int(result.LastSequence))

			if tt.expectedNumEvents == 0 {
				return
			}

			stream, err := js.Stream(ctx, streamName)
			require.NoError(t, err)

			startSeq := result.LastSequence - uint64(result.NumEvents) + 1

			expectedSubject := "test.append.>"
			consumer, err := stream.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{
				OptStartSeq:    startSeq,
				FilterSubjects: []string{expectedSubject},
				DeliverPolicy:  jetstream.DeliverByStartSequencePolicy,
				ReplayPolicy:   jetstream.ReplayInstantPolicy,
			})
			require.NoError(t, err)

			batch, err := consumer.Fetch(tt.expectedNumEvents, jetstream.FetchMaxWait(5*time.Second))
			require.NoError(t, err)

			receivedEvents := []eventstore.Event{}
			for msg := range batch.Messages() {
				data := msg.Data()
				require.NotEmpty(t, data)

				var event cloudevents.Event
				err := event.UnmarshalJSON(data)
				require.NoError(t, err)

				receivedEvents = append(receivedEvents, event)
			}
			require.Equal(t, tt.expectedNumEvents, len(receivedEvents))

			// Compare received events with the expected events.
			for i, event := range receivedEvents {
				expectedEvent := tt.events[i]
				assert.Equal(t, expectedEvent.ID(), event.ID())
				assert.Equal(t, expectedEvent.Type(), event.Type())
				assert.Equal(t, expectedEvent.Source(), event.Source())
				assert.Equal(t, expectedEvent.DataContentType(), event.DataContentType())
				assert.JSONEq(t, string(expectedEvent.Data()), string(event.Data()))
			}
		})
	}
}

func TestEventStoreStream_Pull(t *testing.T) {
	streamName := "TEST_EVENTSTORE_STREAM_PULL"

	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err, "Failed to connect to NATS")
	defer nc.Close()

	type testCase struct {
		name              string
		setup             func(tt *testCase, ctx context.Context, sut eventstore.Stream) error
		batchSize         int
		expectedNumEvents int
		expectedEvents    []eventstore.Event
		expectErr         bool
		expectedErrMsg    string
		expectedErrType   error
	}

	tests := []testCase{
		{
			name: "Pull all events from the stream",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) error {
				events := makeEvents(t, 5, "test.pull")
				_, err := sut.Append(ctx, events)
				if err != nil {
					return fmt.Errorf("failed to append events: %w", err)
				}

				tt.expectedEvents = events
				tt.expectedNumEvents = 5
				return nil
			},
			batchSize:         5,
			expectedNumEvents: 5,
		},
		{
			name: "Pull with empty stream",
			setup: func(tt *testCase, _ context.Context, _ eventstore.Stream) error {
				tt.expectedEvents = []eventstore.Event{}
				tt.expectedNumEvents = 0
				return nil
			},
			batchSize: 5,
		},
	}

	for _, test := range tests {
		tt := test // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, sut := setupTestStream(t, nc, streamName, []string{"test.pull.>"})

			if tt.setup != nil {
				err := tt.setup(&tt, ctx, sut)
				require.NoError(t, err, "Setup function failed")
			}

			pulledEvents, err := sut.Pull(ctx, tt.batchSize)
			if tt.expectErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.IsType(t, tt.expectedErrType, err)
				return
			}

			require.NoError(t, err)

			var idx int
			for _, event := range pulledEvents {
				if idx >=
					len(tt.expectedEvents) {
					t.Fatalf("Received more events than expected: %d > %d", idx, len(tt.expectedEvents))
				}
				expectedEvent := tt.expectedEvents[idx]
				assert.Equal(t, expectedEvent.ID(), event.ID())
				assert.Equal(t, expectedEvent.Type(), event.Type())
				assert.Equal(t, expectedEvent.Source(), event.Source())
				assert.Equal(t, expectedEvent.DataContentType(), event.DataContentType())
				assert.JSONEq(t, string(expectedEvent.Data()), string(event.Data()))
				idx++
			}
			if tt.expectedNumEvents == 0 {
				require.Equal(t, 0, len(pulledEvents), "Expected no events in the stream")
			} else {
				require.Equal(t, tt.expectedNumEvents, len(pulledEvents), "Expected number of pulled events does not match")
			}
		})
	}
}

func TestEventStoreStream_Fetch(t *testing.T) {
	streamName := "TEST_EVENTSTORE_STREAM_FETCH"

	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err, "Failed to connect to NATS")
	defer nc.Close()

	type testCase struct {
		name              string
		setup             func(tt *testCase, ctx context.Context, sut eventstore.Stream) error
		batchSize         int
		subject           string
		expectedNumEvents int
		expectedEvents    []eventstore.Event
		expectErr         bool
		expectedErrMsg    string
		expectedErrType   error
	}

	tests := []testCase{
		{
			name: "Fetch all events from the stream",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) error {
				events := makeEvents(t, 5, "test.fetch")
				_, err := sut.Append(ctx, events)
				if err != nil {
					return fmt.Errorf("failed to append events: %w", err)
				}

				tt.expectedEvents = events
				tt.expectedNumEvents = 5
				return nil
			},
			batchSize:         5,
			expectedNumEvents: 5,
		},
		{
			name: "Fetch events with a subject filter",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) error {
				events1 := makeEvents(t, 5, "test.fetch", "batch-1")
				_, err := sut.Append(ctx, events1)
				if err != nil {
					return fmt.Errorf("failed to append events1: %w", err)
				}

				events2 := makeEvents(t, 5, "test.fetch", "batch-2")
				_, err = sut.Append(ctx, events2)
				if err != nil {
					return fmt.Errorf("failed to append events2: %w", err)
				}

				tt.subject = "test.fetch.event.1"
				tt.expectedEvents = events1
				tt.expectedNumEvents = 5
				return nil
			},
			batchSize: 5,
		},
		{
			name: "Fetch with empty stream",
			setup: func(tt *testCase, _ context.Context, _ eventstore.Stream) error {
				tt.expectedEvents = []eventstore.Event{}
				tt.expectedNumEvents = 0
				return nil
			},
			batchSize: 5,
		},
	}

	for _, test := range tests {
		tt := test // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, sut := setupTestStream(t, nc, streamName, []string{"test.fetch.>"})

			if tt.setup != nil {
				err := tt.setup(&tt, ctx, sut)
				require.NoError(t, err, "Setup function failed")
			}

			opts := []eventstore.FetchOption{
				eventstore.FetchMaxWait(5 * time.Second),
			}
			if tt.subject != "" {
				opts = append(opts, eventstore.FetchSubject(tt.subject))
			}

			fetchedEvents, err := sut.Fetch(ctx, tt.batchSize, opts...)
			if tt.expectErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.IsType(t, tt.expectedErrType, err)
				return
			}

			require.NoError(t, err)

			var idx int
			// Compare fetched events with expected events
			for {
				select {
				case <-ctx.Done():
					t.Fatalf("context cancelled while fetching events: %v", ctx.Err())
				default:
				}

				events, ok := <-fetchedEvents
				require.True(t, ok, "Expected to receive events from the channel")
				if tt.expectedNumEvents == 0 {
					require.Equal(t, 0, len(tt.expectedEvents), "Expected no events in the stream")
					return
				}

				if len(events) == 0 {
					require.Equal(t, tt.batchSize, tt.expectedNumEvents, "Expected no more events")
					return
				}

				if idx == len(tt.expectedEvents) {
					break
				}

				event := events[idx]
				expectedEvent := tt.expectedEvents[idx]
				assert.Equal(t, expectedEvent.ID(), event.ID())
				assert.Equal(t, expectedEvent.Type(), event.Type())
				assert.Equal(t, expectedEvent.Source(), event.Source())
				assert.Equal(t, expectedEvent.DataContentType(), event.DataContentType())
				assert.JSONEq(t, string(expectedEvent.Data()), string(event.Data()))
				idx++
			}
		})
	}
}

func TestEventStoreStream_FetchLast(t *testing.T) {
	streamName := "TEST_EVENTSTORE_STREAM_FETCH_LAST"
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err, "Failed to connect to NATS")
	defer nc.Close()

	type testCase struct {
		name           string
		setup          func(tt *testCase, ctx context.Context, sut eventstore.Stream)
		subject        string
		expectedEvent  eventstore.Event
		expectErr      bool
		expectedErrMsg string
	}

	tests := []testCase{
		{
			name: "Fetch last event from the stream",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events := makeEvents(t, 5, "test.last")
				_, err := sut.Append(ctx, events)
				require.NoError(t, err)
				tt.expectedEvent = events[4]
			},
		},
		{
			name: "Fetch last event with subject filter",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events1 := makeEvents(t, 3, "test.last.1", "batch-1")
				_, err := sut.Append(ctx, events1)
				require.NoError(t, err)

				events2 := makeEvents(t, 1, "test.last.2", "batch-2")
				_, err = sut.Append(ctx, events2)
				require.NoError(t, err)

				tt.subject = "test.last.2.test.event.event_created"
				tt.expectedEvent = events2[0]
			},
		},
		{
			name:           "Fetch last event from empty stream",
			expectErr:      true,
			expectedErrMsg: "event not found",
		},
		{
			name: "Fetch last event with non-existing subject",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events := makeEvents(t, 3, "test.last")
				_, err := sut.Append(ctx, events)
				require.NoError(t, err)
				tt.subject = "non.existing.subject"
			},
			expectErr:      true,
			expectedErrMsg: "event not found",
		},
	}

	for _, test := range tests {
		tt := test // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, sut := setupTestStream(t, nc, streamName, []string{"test.last.>"})

			if tt.setup != nil {
				tt.setup(&tt, ctx, sut)
			}

			opts := []eventstore.FetchOption{}
			if tt.subject != "" {
				opts = append(opts, eventstore.FetchSubject(tt.subject))
			}

			lastEvent, err := sut.FetchLast(ctx, opts...)
			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrMsg)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, lastEvent)

			assert.Equal(t, tt.expectedEvent.ID(), lastEvent.ID())
			assert.Equal(t, tt.expectedEvent.Type(), lastEvent.Type())
			assert.Equal(t, tt.expectedEvent.Source(), lastEvent.Source())
			assert.Equal(t, tt.expectedEvent.DataContentType(), lastEvent.DataContentType())
			assert.JSONEq(t, string(tt.expectedEvent.Data()), string(lastEvent.Data()))
		})
	}
}
