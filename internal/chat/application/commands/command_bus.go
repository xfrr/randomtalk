package chatcommands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/go-cqrsify/cqrs"
	chatdomain "github.com/xfrr/randomtalk/internal/chat/domain"
)

type CommandBus = cqrs.Bus

// InitCommandBus initializes and configures a new command bus instance.
//
// It uses the provided ctx to control request-scoped values, cancellation signals,
// and deadlines during command processing.
func InitCommandBus(
	ctx context.Context,
	csrepo chatdomain.ChatSessionRepository,
	matchRequester chatdomain.MatchRequester,
	logger zerolog.Logger,
) CommandBus {
	cmdbus := cqrs.NewInMemoryBus()

	err := cqrs.Handle(
		ctx,
		cmdbus,
		NewCreateChatSessionCommandHandler(csrepo, matchRequester, logger).Handle,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to register command handler")
	}

	return cmdbus
}

type BaseCommand struct {
	UserAgent      string `json:"-"`
	UserIP         string `json:"-"`
	UserDevice     string `json:"-"`
	UserDeviceType string `json:"-"`
	UserDeviceOS   string `json:"-"`
}
