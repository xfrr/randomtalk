package matchqueries

import (
	"context"

	"github.com/xfrr/go-cqrsify/cqrs"
)

type QueryBus = cqrs.Bus

func InitQueryBus(ctx context.Context) cqrs.Bus {
	qrybus := cqrs.NewInMemoryBus()
	// TODO: add query handlers
	return qrybus
}
