package matchnats

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"github.com/xfrr/randomtalk/internal/shared/eventstore"
	"github.com/xfrr/randomtalk/internal/shared/xnats"
	matchmakingpb "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/matchmaking/v1"
)

var _ matchdomain.NotificationsChannel = (*MatchNotificationsChannel)(nil)

// MatchNotificationsChannel is a channel for subscribing and
// publishing matches to their respective users.
type MatchNotificationsChannel struct {
	streamName string
	jsClient   jetstream.JetStream
	closers    sync.Map
	logger     *zerolog.Logger
}

// Notify sends a notification to the user with the given ID.
func (c *MatchNotificationsChannel) Notify(ctx context.Context, userToNotifyID string, match *matchdomain.Match) error {
	ce, err := createMatchCreatedNotification(*match)
	if err != nil {
		return fmt.Errorf("create match created notification: %w", err)
	}

	body, err := ce.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal match created notification: %w", err)
	}

	natsSubject := buildSubject(
		c.streamName,
		"users",
		userToNotifyID,
		"match_found",
	)

	_, pubErr := c.jsClient.Publish(
		ctx,
		natsSubject,
		body,
	)
	if pubErr != nil {
		return fmt.Errorf("publish notification: %w", pubErr)
	}

	return nil
}

// Subscribe subscribes to the match notifications channel.
func (c *MatchNotificationsChannel) Subscribe(ctx context.Context, userID string) (<-chan string, error) {
	// start pulling from 3 minutes ago
	startTime := time.Now().Add(-3 * time.Minute)

	consumer, err := c.jsClient.CreateConsumer(
		ctx,
		c.streamName,
		jetstream.ConsumerConfig{
			Name:          "match_notifications_" + userID,
			DeliverPolicy: jetstream.DeliverByStartTimePolicy,
			OptStartTime:  &startTime,
			FilterSubjects: []string{
				buildSubject(
					c.streamName,
					"users",
					userID,
					"match_found",
				),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("create match notifications nats consumer: %w", err)
	}

	msgs, err := consumer.Fetch(1)
	if err != nil {
		return nil, fmt.Errorf("fetch messages: %w", err)
	}

	matchIDCh := make(chan string)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msg := <-msgs.Messages()
			if err != nil {
				c.logger.Error().
					Err(err).
					Msg("subscribe to match notifications failed")
				return
			}

			if msg == nil {
				continue
			}

			ce := eventstore.NewEvent()
			if unmarshallErr := ce.UnmarshalJSON(msg.Data()); unmarshallErr != nil {
				c.logger.Error().
					Err(unmarshallErr).
					Msg("unmarshal match created notification")
				continue
			}

			matchIDCh <- ce.Subject()
		}
	}()

	return matchIDCh, nil
}

// CreateMatchNotificationsChannel creates a new MatchNotificationsChannel.
func CreateMatchNotificationsChannel(
	ctx context.Context,
	js jetstream.JetStream,
	logger *zerolog.Logger,
) (*MatchNotificationsChannel, error) {
	streamName := "randomtalk_matchmaking_notifications"

	_, err := createMatchNotificationsStream(ctx, js, streamName)
	if err != nil {
		return nil, fmt.Errorf("create match notifications stream: %w", err)
	}

	return &MatchNotificationsChannel{
		streamName: streamName,
		jsClient:   js,
		logger:     logger,
	}, nil
}

func createMatchCreatedNotification(match matchdomain.Match) (*eventstore.Event, error) {
	ce := eventstore.NewEvent()

	// TODO: inject identity provider
	ce.SetID(uuid.New().String())
	ce.SetType("randomtalk.matchmaking.notifications.match_created_v1")
	ce.SetSource("https://randomtalk.com/matchmaking/notifications")
	ce.SetSpecVersion(cloudevents.CloudEventsVersionV1)
	ce.SetDataContentEncoding(cloudevents.ApplicationJSON)
	ce.SetDataContentType(cloudevents.ApplicationJSON)
	ce.SetSubject(match.ID())
	ce.SetTime(time.Now())
	ce.SetDataSchema(buildNotificationDataSchema("match_created"))

	if err := ce.Context.SetExtension(xnats.SubjectVersionHeaderKey, strconv.Itoa(int(match.AggregateVersion()))); err != nil {
		return nil, fmt.Errorf("set extension: %w", err)
	}

	matchNotification := &matchmakingpb.MatchCreatedNotification{
		MatchId: match.ID(),
		ParticipantIds: []string{
			match.Requester().ID(),
			match.Candidate().ID(),
		},
	}

	if dataErr := ce.SetData(
		string(eventstore.ContentTypeApplicationJSON),
		matchNotification,
	); dataErr != nil {
		return nil, fmt.Errorf("set event data: %w", dataErr)
	}

	return &ce, nil
}
func buildNotificationDataSchema(notificationType string) string {
	return fmt.Sprintf("https://randomtalk.com/schemas/matchmaking/notifications/%s.schema.json", notificationType)
}

func buildSubject(streamName, subjectName string, extra ...string) string {
	streamName = strings.ReplaceAll(streamName, "_", ".")
	subjectName = strings.ReplaceAll(subjectName, "_", ".")

	str := streamName + "." + subjectName
	for _, e := range extra {
		str += "." + e
	}

	return str
}
