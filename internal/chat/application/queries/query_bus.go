package chatqueries

import (
	"context"

	"github.com/xfrr/go-cqrsify/messaging"
)

type QueryBus = messaging.QueryBus

func InitQueryBus(_ context.Context) (QueryBus, func(), error) {
	qrybus := messaging.NewInMemoryQueryBus()
	// TODO: add query handlers
	return qrybus, func() {}, nil
}
