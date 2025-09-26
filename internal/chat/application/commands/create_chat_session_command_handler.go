package chatcommands

import (
	"context"

	"github.com/rs/zerolog"
	chatdomain "github.com/xfrr/randomtalk/internal/chat/domain"
	"github.com/xfrr/randomtalk/internal/chat/infrastructure/auth"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

const CreateChatSessionCommandType = "randomtalk.chat.create_chat_session"

func NewCreateChatSessionCommandHandler(
	chatsessionRepo chatdomain.ChatSessionRepository,
	matchRequester chatdomain.MatchRequester,
	logger zerolog.Logger,
) CreateChatSessionCommandHandler {
	return CreateChatSessionCommandHandler{
		logger:          logger,
		matchRequester:  matchRequester,
		chatSessionRepo: chatsessionRepo,
	}
}

type CreateChatSessionCommandHandler struct {
	logger          zerolog.Logger
	matchRequester  chatdomain.MatchRequester
	chatSessionRepo chatdomain.ChatSessionRepository
}

func (h CreateChatSessionCommandHandler) Handle(ctx context.Context, cmd CreateChatSessionCommand) error {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return ErrMissingUserIDFromContext
	}

	h.logger.Debug().
		Str("user_id", userID).
		Str("user_nickname", cmd.UserNickname).
		Int32("user_age", cmd.UserAge).
		Msg("an user requested a new random chat session")

	if exists, err := h.chatSessionRepo.Exists(ctx, userID); err != nil {
		return err
	} else if exists {
		return chatdomain.
			ErrChatSessionAlreadyExists.
			WithAggregateID(userID).
			WithAggregateName(chatdomain.AggregateName)
	}

	user, err := chatdomain.NewUser(
		chatdomain.ID(userID),
		cmd.UserNickname,
		cmd.UserAge,
		gender.Parse(cmd.UserGender),
		matchmaking.DefaultPreferences().
			WithMinAge(cmd.UserMatchPreferenceMinAge).
			WithMaxAge(cmd.UserMatchPreferenceMaxAge).
			WithGender(gender.Parse(cmd.UserMatchPreferenceGender)).
			WithInterests(cmd.UserMatchPreferenceInterests),
	)
	if err != nil {
		return err
	}

	cs, err := chatdomain.NewChatSession(user.ID(), user)
	if err != nil {
		return err
	}

	if err = h.chatSessionRepo.Save(ctx, cs); err != nil {
		return err
	}

	if err = h.matchRequester.RequestMatch(ctx, cs); err != nil {
		return err
	}

	return nil
}
