package matchnats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// MatchNotificationsStream is a stream for publishing Matchmaking notifications.
type MatchNotificationsStream struct {
	name string
	js   jetstream.JetStream
}

func (s *MatchNotificationsStream) Name() string {
	return s.name
}

func createMatchNotificationsStream(ctx context.Context, js jetstream.JetStream, streamName string) (*MatchNotificationsStream, error) {
	// create the stream if it doesn't exist
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      streamName,
		Retention: jetstream.InterestPolicy,
		MaxAge:    5 * time.Minute,
		Subjects:  []string{buildSubject(streamName, ">")},
	})
	if err != nil {
		return nil, fmt.Errorf("create nats stream: %w", err)
	}

	return &MatchNotificationsStream{
		name: streamName,
		js:   js,
	}, nil
}
