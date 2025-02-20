//go:build integration
// +build integration

package matchnats_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	matchsessionreponats "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/nats"
	eventstore "github.com/xfrr/randomtalk/pkg/eventstore"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStream(
	t *testing.T,
	streamName string,
	subjects []string,
) (context.Context, jetstream.JetStream, eventstore.Stream) {
	t.Helper()

	ctx := context.Background()
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err, "Failed to connect to NATS")

	js, err := jetstream.New(nc)
	require.NoError(t, err, "Failed to create JetStream context")

	sut, err := matchsessionreponats.CreateStream(ctx, js, matchsessionreponats.NewStreamConfig("test-stream"))
	require.NoError(t, err)

	// Clean up resources after the test.
	t.Cleanup(func() {
		_ = js.DeleteConsumer(ctx, streamName, streamName)
		_ = js.DeleteStream(ctx, streamName)
		nc.Close()
	})

	return ctx, js, sut
}

func makeEvent(t *testing.T, id string) eventstore.Event {
	t.Helper()
	event := cloudevents.NewEvent()
	event.SetID(fmt.Sprintf("test-event-%s", id))
	event.SetType("event_created")
	event.SetSource("test.eventstore")
	event.SetSubject("test.event")
	event.SetDataContentType("application/json")
	err := event.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"event": "test",
	})
	require.NoError(t, err)

	return event
}

func makeEvents(t *testing.T, n int, prefix ...string) []eventstore.Event {
	t.Helper()
	events := make([]eventstore.Event, n)
	for i := 0; i < n; i++ {
		eventID := fmt.Sprintf("%s-%d", strings.Join(prefix, ""), i)
		events[i] = makeEvent(t, eventID)
	}
	return events
}

func TestEventStoreStream_Append(t *testing.T) {
	streamName := "TEST_EVENTSTORE_STREAM_APPEND"
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
			name:              "Append a list of events to an empty stream",
			events:            makeEvents(t, 5),
			expectedNumEvents: 5,
			expectedLastSeq:   5,
			expectErr:         false,
		},
		{
			name:              "Append a list of events to a non-empty stream",
			events:            makeEvents(t, 5),
			expectedNumEvents: 10,
			expectedLastSeq:   10,
			expectErr:         false,
		},
		{
			name:              "Append with empty events list",
			events:            []eventstore.Event{},
			expectedNumEvents: 0,
			expectedLastSeq:   0,
			expectErr:         false,
		},
	}

	for _, test := range tests {
		tt := test // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, js, sut := setupTestStream(t, streamName, []string{"test.append.>"})

			result, err := sut.Append(ctx, tt.events)
			if tt.expectErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErrMsg)
				require.IsType(t, tt.expectedErrType, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedNumEvents, result.NumEvents)
			require.Equal(t, tt.expectedLastSeq, result.LastSequence)

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

func TestEventStoreStream_Fetch(t *testing.T) {
	streamName := "TEST_EVENTSTORE_STREAM_FETCH"

	type testCase struct {
		name              string
		setup             func(tt *testCase, ctx context.Context, sut eventstore.Stream)
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
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events := makeEvents(t, 5)
				_, err := sut.Append(ctx, events)
				require.NoError(t, err)
				tt.expectedEvents = events
			},
			batchSize:         5,
			expectedNumEvents: 5,
		},
		{
			name: "Fetch events with a subject filter",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events1 := makeEvents(t, 5, "batch-1")
				_, err := sut.Append(ctx, events1)
				require.NoError(t, err)

				events2 := makeEvents(t, 5, "batch-2")
				_, err = sut.Append(ctx, events2)
				require.NoError(t, err)

				tt.subject = "test.fetch.event.1"
				tt.expectedEvents = events1
				tt.expectedNumEvents = 5
			},
			batchSize: 5,
		},
		{
			name: "Fetch with batch size larger than available events",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events := makeEvents(t, 3)
				_, err := sut.Append(ctx, events)
				require.NoError(t, err)
				tt.expectedEvents = events
				tt.expectedNumEvents = 3
			},
			batchSize: 5,
		},
		{
			name: "Fetch with empty stream",
			setup: func(tt *testCase, _ context.Context, _ eventstore.Stream) {
				tt.expectedEvents = []eventstore.Event{}
				tt.expectedNumEvents = 0
			},
			batchSize: 5,
		},
	}

	for _, test := range tests {
		tt := test // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, sut := setupTestStream(t, streamName, []string{"test.fetch.>"})

			if tt.setup != nil {
				tt.setup(&tt, ctx, sut)
			}

			opts := []eventstore.FetchOption{}
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
			require.Len(t, fetchedEvents, tt.expectedNumEvents)

			var idx int
			// Compare fetched events with expected events
			for events := range fetchedEvents {
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
	type testCase struct {
		name            string
		setup           func(tt *testCase, ctx context.Context, sut eventstore.Stream)
		subject         string
		expectedEvent   eventstore.Event
		expectErr       bool
		expectedErrMsg  string
		expectedErrType error
	}

	tests := []testCase{
		{
			name: "Fetch last event from the stream",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events := makeEvents(t, 5)
				_, err := sut.Append(ctx, events)
				require.NoError(t, err)
				tt.expectedEvent = events[4]
			},
		},
		{
			name: "Fetch last event with subject filter",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events1 := makeEvents(t, 3, "batch-1")
				_, err := sut.Append(ctx, events1)
				require.NoError(t, err)

				events2 := makeEvents(t, 2, "batch-2")
				_, err = sut.Append(ctx, events2)
				require.NoError(t, err)

				tt.subject = "test.last.event.1"
				tt.expectedEvent = events1[2]
			},
		},
		{
			name:            "Fetch last event from empty stream",
			expectErr:       true,
			expectedErrMsg:  "bad request",
			expectedErrType: jetstream.ErrBadRequest,
		},
		{
			name: "Fetch last event with non-existing subject",
			setup: func(tt *testCase, ctx context.Context, sut eventstore.Stream) {
				events := makeEvents(t, 3)
				_, err := sut.Append(ctx, events)
				require.NoError(t, err)
				tt.subject = "non.existing.subject"
			},
			expectErr:      true,
			expectedErrMsg: "message not found",
		},
	}

	for _, test := range tests {
		tt := test // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, sut := setupTestStream(t, streamName, []string{"test.last.>"})

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

func TestEventStoreStream_AppendConcurrentSameSubject(t *testing.T) {
	streamName := "TEST_EVENTSTORE_STREAM_APPEND_CONCURRENT_SAME_SUBJECT"
	ctx, js, sut := setupTestStream(t, streamName, []string{"test.concurrent.same.subject"})

	events := makeEvents(t, 5)

	go func() {
		_, err := sut.Append(ctx, events)
		require.NoError(t, err)
	}()

	// consume the events
	stream, err := js.Stream(ctx, streamName)
	require.NoError(t, err)

	consumer, err := stream.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{
		FilterSubjects: []string{"test.concurrent.same.subject"},
		DeliverPolicy:  jetstream.DeliverAllPolicy,
		ReplayPolicy:   jetstream.ReplayInstantPolicy,
	})
	require.NoError(t, err)

	batch, err := consumer.Fetch(5, jetstream.FetchMaxWait(5*time.Second))
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

	require.Len(t, receivedEvents, 5)
	for i, event := range receivedEvents {
		expectedEvent := events[i]
		assert.Equal(t, expectedEvent.ID(), event.ID())
		assert.Equal(t, expectedEvent.Type(), event.Type())
		assert.Equal(t, expectedEvent.Source(), event.Source())
		assert.Equal(t, expectedEvent.DataContentType(), event.DataContentType())
		assert.JSONEq(t, string(expectedEvent.Data()), string(event.Data()))
	}
}
