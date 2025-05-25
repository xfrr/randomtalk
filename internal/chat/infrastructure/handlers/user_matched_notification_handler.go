package chatinfrahandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/xfrr/randomtalk/internal/shared/messaging"
	"github.com/xfrr/randomtalk/internal/shared/semantic"

	chatcommands "github.com/xfrr/randomtalk/internal/chat/application/commands"
)

type UserMatchedNotification struct {
	MatchID string `json:"match_id"`
	UserIDs string `json:"user_ids"`
}

type UserMatchedNotificationHandler struct {
	cmdbus chatcommands.CommandBus
	logger *zerolog.Logger
}

func NewUserMatchedNotificationHandler(
	cmdbus chatcommands.CommandBus,
	logger *zerolog.Logger,
) *UserMatchedNotificationHandler {
	return &UserMatchedNotificationHandler{
		cmdbus: cmdbus,
		logger: logger,
	}
}

func (h *UserMatchedNotificationHandler) Handle(ctx context.Context, evt *messaging.Event) error {
	h.logger.Debug().
		Str(semantic.MessageIDKey, evt.ID()).
		Str(semantic.MessageTypeKey, evt.Type()).
		Msg("user matched notification received")

	// parse event payload
	eventPayload := new(UserMatchedNotification)
	err := json.Unmarshal(evt.Data(), eventPayload)
	if err != nil {
		// reject event
		evt.Reject()
		return fmt.Errorf("unmarshal user matched notification: %w", err)
	}

	// TODO: dispatch command

	// ack event
	evt.Ack()
	return nil
}
