package matchcommands

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/xfrr/go-cqrsify/cqrs"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
)

type CommandBus = cqrs.Bus

func InitCommandBus(
	ctx context.Context,
	matchmakingService matchdomain.MatchmakingProcessor,
) cqrs.Bus {
	cmdbus := cqrs.NewInMemoryBus()

	err := cqrs.Handle(
		ctx,
		cmdbus,
		NewMatchmakingCommandHandler(matchmakingService).ProcessMatchUserWithPreferencesCommand,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register match user with preferences command handler")
	}

	return cmdbus
}
