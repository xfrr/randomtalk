package matchmakinghandlers

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
	"github.com/xfrr/randomtalk/internal/shared/messaging"
	chatpbv1 "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/chat/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserMatchRequestedNotificationHandler struct {
	logger               *zerolog.Logger
	matchmakingProcessor matchdomain.MatchmakingProcessor
}

func NewUserMatchRequestedEventHandler(
	matchmakingService matchdomain.MatchmakingProcessor,
	logger *zerolog.Logger,
) *UserMatchRequestedNotificationHandler {
	return &UserMatchRequestedNotificationHandler{
		logger:               logger,
		matchmakingProcessor: matchmakingService,
	}
}

func (h *UserMatchRequestedNotificationHandler) Handle(ctx context.Context, msg *messaging.Event) error {
	h.logger.Debug().
		Str("messaging_event_id", msg.ID()).
		Str("messaging_event_type", msg.Type()).
		Msg("user match requested notification received")

	notification := new(chatpbv1.UserMatchRequestedNotification)
	err := protojson.Unmarshal(msg.Data(), notification)
	if err != nil {
		// discard message
		msg.Nack()
		return fmt.Errorf("unmarshal user match requested notification: %w", err)
	}

	// create user from notification
	user := matchdomain.NewUser(
		notification.GetUserAttributes().GetId(),
		notification.GetUserAttributes().GetAge(),
		toGender(notification.GetUserAttributes().GetGender()),
		matchmaking.DefaultPreferences().
			WithMinAge(notification.GetUserPreferences().GetMinAge()).
			WithMaxAge(notification.GetUserPreferences().GetMaxAge()).
			WithGender(toGender(notification.GetUserPreferences().GetGender())).
			WithInterests(notification.GetUserPreferences().GetInterests()),
	)

	// attempt to match user with preferences
	err = h.matchmakingProcessor.ProcessMatchRequest(ctx, *user)
	if err != nil {
		// nack msg to retry
		msg.Nack()
		return fmt.Errorf("attempt match with preferences: %w", err)
	}

	// ack msg
	msg.Ack()
	return nil
}

func toGender(g chatpbv1.Gender) gender.Gender {
	switch g {
	case chatpbv1.Gender_GENDER_FEMALE:
		return gender.Female
	case chatpbv1.Gender_GENDER_MALE:
		return gender.Male
	default:
		return gender.Unspecified
	}
}
