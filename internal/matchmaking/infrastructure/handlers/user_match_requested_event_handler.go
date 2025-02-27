package matchmakinginfrahandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/location"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
	"github.com/xfrr/randomtalk/internal/shared/messaging"
)

type UserMatchRequestedEvent struct {
	UserID       string                        `json:"user_id"`
	UserName     string                        `json:"user_name"`
	UserAge      int                           `json:"user_age"`
	UserGender   gender.Gender                 `json:"user_gender"`
	UserLocation *location.Location            `json:"user_location,omitempty"`
	Preferences  *matchmaking.MatchPreferences `json:"user_preferences,omitempty"`
}

type UserMatchRequestedEventHandler struct {
	logger               *zerolog.Logger
	matchmakingProcessor matchdomain.MatchmakingProcessor
}

func NewUserMatchRequestedEventHandler(
	matchmakingService matchdomain.MatchmakingProcessor,
	logger *zerolog.Logger,
) *UserMatchRequestedEventHandler {
	return &UserMatchRequestedEventHandler{
		logger:               logger,
		matchmakingProcessor: matchmakingService,
	}
}

func (h *UserMatchRequestedEventHandler) Handle(ctx context.Context, evt *messaging.Event) error {
	h.logger.Debug().
		Str("messaging_event_id", evt.ID()).
		Str("messaging_event_type", evt.Type()).
		Msg("user match requested event received")

	// parse event payload
	eventPayload := new(UserMatchRequestedEvent)
	err := json.Unmarshal(evt.Data(), eventPayload)
	if err != nil {
		// reject event
		evt.Reject()
		return fmt.Errorf("unmarshal user: %w", err)
	}

	// create user from event
	user := matchdomain.NewUser(
		eventPayload.UserID,
		eventPayload.UserAge,
		eventPayload.UserGender,
		eventPayload.Preferences,
	)

	// attempt to match user with preferences
	err = h.matchmakingProcessor.ProcessMatchRequest(ctx, user)
	if err != nil {
		// nack event to retry
		evt.Nack()
		return fmt.Errorf("attempt match with preferences: %w", err)
	}

	// ack event
	evt.Ack()
	return nil
}
