//go:build integration
// +build integration

package randomtalk_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	natsURL          = nats.DefaultURL
	chatSubjectBase  = "randomtalk.chat.notifications.sessions"
	matchmakingTopic = "randomtalk.matchmaking.matches.>"
	numberOfUsers    = 10
	testTimeout      = 5 * time.Minute // Timeout for the entire test
)

var (
	expectedMatches = numberOfUsers / 2 // With 1000 users, we expect 500 matches
)

type NATSClient struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func TestRandomtalkIntegration(t *testing.T) {
	setupLogger()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	client := mustCreateNATSClient(t, natsURL)
	defer client.Close()

	// 1) Start publishing chat notifications in a separate goroutine.
	//    We do not tie this to the main wait group; it simply fires and forgets.
	go publishChatNotifications(ctx, t, client)

	// 2) Channel to receive the match results from the matching listener.
	matchResults := make(chan matchResult, 1)

	// 3) Start listening for matches in another goroutine.
	go func() {
		matchesFound, duplicatedMatches := listenForMatches(ctx, t, client)
		matchResults <- matchResult{
			matchesFound:      matchesFound,
			duplicatedMatches: duplicatedMatches,
		}
		close(matchResults)
	}()

	// 4) Wait for the result (or test context expiration).
	select {
	case res, ok := <-matchResults:
		require.True(t, ok, "Expected match results channel to be open")
		require.Zero(t, res.duplicatedMatches, "Unexpected duplicate matches found")
		assert.Equal(t, expectedMatches, res.matchesFound, "Unexpected number of matches found")
	case <-ctx.Done():
		t.Fatal("Test timed out before receiving match results")
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // Verbose for integration tests
}

func mustCreateNATSClient(t *testing.T, url string) *NATSClient {
	client, err := NewNATSClient(url)
	require.NoError(t, err, "Failed to initialize NATS client")
	return client
}

func NewNATSClient(url string) (*NATSClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &NATSClient{nc: nc, js: js}, nil
}

func (n *NATSClient) Close() {
	// Drain returns an error if called multiple times or the connection is closed,
	// so you may ignore the error or log it if you wish.
	_ = n.nc.Drain()
	n.nc.Close()
}

func publishChatNotifications(ctx context.Context, t *testing.T, client *NATSClient) {
	var wg sync.WaitGroup
	wg.Add(numberOfUsers)

	for i := 1; i <= numberOfUsers; i++ {
		userID := i
		go func() {
			defer wg.Done()
			publishUserNotification(ctx, t, client, userID)
		}()
	}

	wg.Wait()

	// Flush to ensure all messages are actually sent before we exit this goroutine.
	// If flush fails, it's likely a bigger issue, so we fail the test immediately.
	require.NoError(t, client.nc.Flush(), "NATS flush failed after publishing chat notifications")

	log.Info().Msg("All chat notifications published")
}

func publishUserNotification(ctx context.Context, t *testing.T, client *NATSClient, userID int) {
	event := createChatEvent(userID)
	body, err := json.Marshal(event)
	require.NoError(t, err, "Failed to marshal CloudEvent")

	subject := fmt.Sprintf("%s.%d.connected", chatSubjectBase, userID)
	_, err = client.js.Publish(ctx, subject, body, jetstream.WithMsgID(event.ID()))
	require.NoErrorf(t, err, "Failed to publish message for user %d", userID)
}

func createChatEvent(userID int) cloudevents.Event {
	event := cloudevents.New()
	event.SetID(fmt.Sprintf("chat_session_started_%d", userID))
	event.SetType("chat_session_started")
	event.SetSource("https://randomtalk.com/chat/notifications")
	event.SetSubject(fmt.Sprintf("user.%d", userID))
	event.SetTime(time.Now())
	event.SetDataContentType(cloudevents.ApplicationJSON)

	data := map[string]interface{}{
		"user_attributes": map[string]interface{}{
			"id":     strconv.Itoa(userID),
			"age":    20,
			"gender": "GENDER_MALE",
		},
		"user_preferences": map[string]interface{}{
			"min_age":   20,
			"max_age":   30,
			"interests": []string{"music", "sports"},
		},
	}
	event.SetData(data)
	return event
}

// matchResult is passed through a channel to the main goroutine
type matchResult struct {
	matchesFound      int
	duplicatedMatches int
}

func listenForMatches(ctx context.Context, t *testing.T, client *NATSClient) (int, int) {
	consumer, err := client.js.OrderedConsumer(ctx,
		"randomtalk_matchmaking_match_events", // Stream name
		jetstream.OrderedConsumerConfig{
			FilterSubjects:    []string{matchmakingTopic},
			InactiveThreshold: 2 * time.Minute,
		})
	require.NoError(t, err, "failed to create ordered consumer")

	var (
		matchesFound      int
		duplicatedMatches int
	)
	matchesSeen := sync.Map{}

	for {
		// If we've already got all the matches, exit early for faster tests
		if matchesFound >= expectedMatches {
			log.Info().Msgf("Reached expected match count: %d", matchesFound)
			break
		}

		select {
		case <-ctx.Done():
			log.Warn().Msg("Stopped listening for matches due to context timeout/cancellation")
			break
		default:
			// Try to fetch a batch of messages
			batch, fetchErr := consumer.Fetch(100, jetstream.FetchMaxWait(2*time.Second))
			if fetchErr != nil {
				if fetchErr == nats.ErrTimeout {
					// No new messages in the last 2s; continue listening.
					continue
				}
				t.Fatalf("Failed to fetch messages: %v", fetchErr)
			}

			for msg := range batch.Messages() {
				match, procErr := processMatchMessage(msg)
				if procErr != nil {
					t.Logf("Failed to process match message: %v", procErr)
					continue
				}

				// Check if the requester or candidate was seen before
				if markAsSeen(&matchesSeen, match.RequesterID) || markAsSeen(&matchesSeen, match.CandidateID) {
					duplicatedMatches++
				}

				matchesFound++
				require.NoError(t, msg.Ack())
			}
		}
		// No 'continue' or 'break' here means we re-check the condition at the top of the loop
	}

	return matchesFound, duplicatedMatches
}

// markAsSeen returns true if the userID was already stored in the map.
// We store an empty struct to reduce memory overhead.
func markAsSeen(matchesUniqueMap *sync.Map, userID string) bool {
	_, loaded := matchesUniqueMap.LoadOrStore(userID, struct{}{})
	return loaded
}

type matchPayload struct {
	RequesterID string `json:"match_user_requester_id"`
	CandidateID string `json:"match_user_matched_id"`
}

func processMatchMessage(msg jetstream.Msg) (*matchPayload, error) {
	var evt cloudevents.Event
	if err := json.Unmarshal(msg.Data(), &evt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CloudEvent: %w", err)
	}

	var p matchPayload
	if err := evt.DataAs(&p); err != nil {
		return nil, fmt.Errorf("failed to parse event data: %w", err)
	}

	return &p, nil
}
