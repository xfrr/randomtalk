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
	subjectTemplate  = "randomtalk.notifications.chat.users"
	matchmakingTopic = "randomtalk.matchmaking.matches.>"
	numberOfUsers    = 500
	expectedMatches  = 250
	testTimeout      = 10 * time.Minute
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		publishChatNotifications(ctx, t, client)
	}()

	matchResults := make(chan struct {
		matchesFound      int
		duplicatedMatches int
	}, 1)

	go func() {
		matchesFound, duplicatedMatches := listenForMatches(ctx, t, client)
		matchResults <- struct {
			matchesFound      int
			duplicatedMatches int
		}{matchesFound, duplicatedMatches}
		close(matchResults)
	}()

	wg.Wait() // Ensure all messages are published before assertions
	select {
	case res := <-matchResults:
		require.Zero(t, res.duplicatedMatches, "unexpected duplicate matches found")
		assert.Equal(t, expectedMatches, res.matchesFound, "unexpected number of matches found")
	case <-ctx.Done():
		t.Fatal("Test timed out before receiving match results")
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func mustCreateNATSClient(t *testing.T, url string) *NATSClient {
	nc, err := nats.Connect(url)
	require.NoError(t, err, "failed to connect to NATS")

	js, err := jetstream.New(nc)
	require.NoError(t, err, "failed to create JetStream context")

	return &NATSClient{nc: nc, js: js}
}

func (n *NATSClient) Close() {
	_ = n.nc.Drain()
	n.nc.Close()
}

func publishChatNotifications(ctx context.Context, t *testing.T, client *NATSClient) {
	var wg sync.WaitGroup
	wg.Add(numberOfUsers)

	for i := 1; i <= numberOfUsers; i++ {
		go func(userID int) {
			defer wg.Done()
			publishUserNotification(ctx, t, client, userID)
		}(i)
	}

	waitForCompletion(ctx, &wg, t, "publishChatNotifications")
}

func publishUserNotification(ctx context.Context, t *testing.T, client *NATSClient, userID int) {
	event := createChatEvent(userID)
	body, err := json.Marshal(event)
	if err != nil {
		t.Logf("Failed to marshal CloudEvent: %v", err)
		return
	}

	subject := fmt.Sprintf("%s.%d.connected", subjectTemplate, userID)
	_, err = client.js.Publish(ctx, subject, body, jetstream.WithMsgID(event.ID()))
	if err != nil {
		t.Logf("Failed to publish message for user %d: %v", userID, err)
	}
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
		"user_id":     strconv.Itoa(userID),
		"user_age":    20,
		"user_gender": "GENDER_MALE",
		"user_preferences": map[string]interface{}{
			"min_age":   20,
			"max_age":   30,
			"interests": []string{"music", "sports"},
		},
	}
	event.SetData(data)
	return event
}

func listenForMatches(ctx context.Context, t *testing.T, client *NATSClient) (int, int) {
	consumer, err := client.js.OrderedConsumer(ctx, "randomtalk_matchmaking_match_events",
		jetstream.OrderedConsumerConfig{FilterSubjects: []string{matchmakingTopic}, InactiveThreshold: 2 * time.Minute})
	require.NoError(t, err, "failed to create ordered consumer")

	var matchesFound, duplicatedMatches int
	matchesUniqueMap := sync.Map{}
	msgChan := make(chan jetstream.Msg, 1000)

	go func() {
		consumer.Consume(func(msg jetstream.Msg) {
			msgChan <- msg
		})
	}()

	for {
		select {
		case <-ctx.Done():
			return matchesFound, duplicatedMatches
		case msg := <-msgChan:
			match, err := processMatchMessage(msg)
			if err != nil {
				t.Logf("Failed to process match message: %v", err)
				continue
			}

			if isDuplicateMatch(&matchesUniqueMap, match.RequesterID) || isDuplicateMatch(&matchesUniqueMap, match.CandidateID) {
				duplicatedMatches++
				t.Logf("Duplicate match found: %+v", match)
			}
			matchesFound++
		}
	}
}

func isDuplicateMatch(matchesUniqueMap *sync.Map, userID string) bool {
	_, loaded := matchesUniqueMap.LoadOrStore(userID, struct{}{})
	return loaded
}

type payload struct {
	RequesterID string `json:"match_user_requester_id"`
	CandidateID string `json:"match_user_matched_id"`
}

func processMatchMessage(msg jetstream.Msg) (*payload, error) {
	var event cloudevents.Event
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CloudEvent: %w", err)
	}

	var p *payload
	if err := event.DataAs(&p); err != nil {
		return nil, fmt.Errorf("failed to parse event data: %w", err)
	}

	if err := msg.Ack(); err != nil {
		return nil, fmt.Errorf("failed to acknowledge message: %w", err)
	}

	return p, nil
}

func waitForCompletion(ctx context.Context, wg *sync.WaitGroup, t *testing.T, name string) {
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		t.Logf("%s canceled or timed out: %v", name, ctx.Err())
	case <-doneCh:
		t.Logf("%s completed.", name)
	}
}
