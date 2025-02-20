package matchnats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// CreateMatchmakingUserMatchRequestsStream creates a JetStream stream for matchmaking user match requests.
func CreateMatchmakingUserMatchRequestsStream(ctx context.Context, js jetstream.JetStream) error {
	// create user match requested stream
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:       "randomtalk_matchmaking_user_match_requests",
		Retention:  jetstream.WorkQueuePolicy, // exactly-once delivery
		Subjects:   []string{"randomtalk.matchmaking.user_match_requests.>"},
		MaxAge:     5 * time.Minute,
		MaxMsgSize: 1024, // 1KB
	})
	if err != nil {
		return err
	}

	return nil
}
