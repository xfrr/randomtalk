package chatnats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	chatdomain "github.com/xfrr/randomtalk/internal/chat/domain"
	"github.com/xfrr/randomtalk/internal/shared/eventstore"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	chatpbv1 "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/chat/v1"
)

const (
	// EventTypeUserMatchRequested is the CloudEvent type for user match request events.
	EventTypeUserMatchRequested = "com.randomtalk.chat.notifications.user_match_requested"
	// EventSource identifies the source of the event.
	EventSource = "/chat"
	// maxRetries defines the number of retry attempts for publishing.
	maxRetries = 3
	// retryDelay is the base delay duration between retries.
	retryDelay = 500 * time.Millisecond
)

var _ chatdomain.MatchRequester = (*MatchRequester)(nil)

// MatchRequester is responsible for publishing user match request events via NATS JetStream.
// It uses CloudEvents to ensure interoperability and standardization.
type MatchRequester struct {
	streamName string
	js         jetstream.JetStream
}

// It creates a CloudEvent based on the provided ChatSession, marshals it to JSON,
// and publishes it to the NATS JetStream stream.
func (m *MatchRequester) RequestMatch(ctx context.Context, cs *chatdomain.ChatSession) error {
	eventID := uuid.New().String()

	ce := eventstore.NewEvent()
	ce.SetID(eventID)
	ce.SetType(EventTypeUserMatchRequested)
	ce.SetSource(chatdomain.EventSourceName)
	ce.SetSubject(strings.Join([]string{chatSessionsStreamSuffix, cs.AggregateID()}, "."))
	ce.SetTime(time.Now().UTC())
	ce.SetDataSchema("schemas.randomtalk.com/chat/notifications/user_match_requested/1.0")

	notif := &chatpbv1.UserMatchRequestedNotification{
		NotificationId: eventID,
		ChatSessionId:  cs.AggregateID(),
		UserAttributes: &chatpbv1.UserAttributes{
			Id:     cs.User().ID().String(),
			Age:    cs.User().Age(),
			Gender: toProtoGender(cs.User().Gender()),
		},
		UserPreferences: &chatpbv1.UserPreferences{
			MinAge:    cs.User().MatchPreferences().MinAge,
			MaxAge:    cs.User().MatchPreferences().MaxAge,
			Gender:    toProtoGender(cs.User().MatchPreferences().Gender),
			Interests: cs.User().MatchPreferences().Interests,
		},
	}

	if dataErr := ce.SetData(string(eventstore.ContentTypeApplicationJSON), notif); dataErr != nil {
		return fmt.Errorf("set event data: %w", dataErr)
	}

	body, err := ce.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal user match request cloudevent notification: %w", err)
	}

	subject := strings.Join([]string{chatdomain.EventSourceName, "notifications", cs.AggregateID(), "user_match_requested"}, ".")

	msg := nats.Msg{
		Subject: subject,
		Data:    body,
	}

	_, err = m.js.PublishMsg(ctx, &msg,
		jetstream.WithExpectStream(m.streamName),
		jetstream.WithMsgID(eventID),
		jetstream.WithRetryAttempts(maxRetries),
		jetstream.WithRetryWait(retryDelay),
	)
	if err != nil {
		return fmt.Errorf("publish user match request event: %w", err)
	}

	return nil
}

func NewMatchRequester(streamName string, js jetstream.JetStream) *MatchRequester {
	return &MatchRequester{
		streamName: streamName,
		js:         js,
	}
}

func toProtoGender(g gender.Gender) chatpbv1.Gender {
	switch g {
	case gender.Female:
		return chatpbv1.Gender_GENDER_FEMALE
	case gender.Male:
		return chatpbv1.Gender_GENDER_MALE
	default:
		return chatpbv1.Gender_GENDER_UNSPECIFIED
	}
}
