package chatcommands

import (
	"context"

	"github.com/xfrr/go-cqrsify/cqrs"
)

type CommandBus = cqrs.Bus

func InitCommandBus(ctx context.Context) cqrs.Bus {
	cmdbus := cqrs.NewInMemoryBus()

	// TODO: add command handlers
	return cmdbus
}
