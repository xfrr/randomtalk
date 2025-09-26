package matchcommands

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/xfrr/go-cqrsify/messaging"

	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
)

type CommandBus = messaging.CommandBus

func InitCommandBus(
	ctx context.Context,
	matchmakingService matchdomain.MatchmakingProcessor,
) (CommandBus, func()) {
	cmdbus := messaging.NewInMemoryCommandBus()

	unsub, err := messaging.SubscribeCommand(
		ctx,
		cmdbus,
		MatchUserWithPreferencesCommandType,
		NewMatchmakingCommandHandler(matchmakingService),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register match user with preferences command handler")
	}

	closer := func() {
		unsub()
	}

	return cmdbus, closer
}
