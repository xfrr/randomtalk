package chatcommands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/go-cqrsify/messaging"
	chatdomain "github.com/xfrr/randomtalk/internal/chat/domain"
)

type CommandBus = messaging.CommandBus

// InitCommandBus initializes and configures a new command bus instance.
//
// It uses the provided ctx to control request-scoped values, cancellation signals,
// and deadlines during command processing.
func InitCommandBus(
	ctx context.Context,
	csrepo chatdomain.ChatSessionRepository,
	matchRequester chatdomain.MatchRequester,
	logger zerolog.Logger,
) (CommandBus, func(), error) {
	cmdbus := messaging.NewInMemoryCommandBus()

	unsubCreateChatSessionCmd, err := messaging.SubscribeCommand(
		ctx,
		cmdbus,
		CreateChatSessionCommandType,
		NewCreateChatSessionCommandHandler(csrepo, matchRequester, logger),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to register command handler")
	}

	closer := func() {
		unsubCreateChatSessionCmd()
	}

	return cmdbus, closer, nil
}

type CommandInfo struct {
	UserAgent      string `json:"-"`
	UserIP         string `json:"-"`
	UserDevice     string `json:"-"`
	UserDeviceType string `json:"-"`
	UserDeviceOS   string `json:"-"`
}
