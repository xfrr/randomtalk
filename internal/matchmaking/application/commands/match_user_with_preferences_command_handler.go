package matchcommands

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
)

var (
	ErrUnableToMatch = errors.New("unable to match user with the given preferences")
)

type MatchmakingCommandHandler struct {
	matchmakingService matchdomain.MatchmakingProcessor
}

func NewMatchmakingCommandHandler(ms matchdomain.MatchmakingProcessor) *MatchmakingCommandHandler {
	return &MatchmakingCommandHandler{
		matchmakingService: ms,
	}
}

func (h *MatchmakingCommandHandler) ProcessMatchUserWithPreferencesCommand(
	ctx context.Context,
	cmd MatchUserWithPreferencesCommand,
) (interface{}, error) {
	log.Debug().
		Str("user_id", cmd.UserID).
		Int("user_age", cmd.UserAge).
		Any("user_preferences", cmd.UserMatchPreferences).
		Msg("match user command received")

	requesterUser := matchdomain.NewUser(
		cmd.UserID,
		cmd.UserAge,
		cmd.UserGender,
		&cmd.UserMatchPreferences,
	)

	err := h.matchmakingService.ProcessMatchRequest(ctx, requesterUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find best match: %w", err)
	}

	// TODO: subscribe to the match notifications channel to receive the match

	return &MatchUserWithPreferencesResponse{
		// MatchID: match.ID(),
	}, nil
}
