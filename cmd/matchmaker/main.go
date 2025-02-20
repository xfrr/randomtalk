package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/xfrr/randomtalk/internal/matchmaking"

	"github.com/xfrr/randomtalk/pkg/env"
)

var (
	// ServiceVersion is configurable using go build -ldflags "-X main.ServiceVersion=..."
	// or by setting the SERVICE_VERSION environment variable.
	ServiceVersion = env.GetWithDefault("SERVICE_VERSION", "development")
)

func main() {
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	svc := matchmaking.MustInitService(ctx, matchmaking.WithVersion(ServiceVersion))
	defer svc.Shutdown()

	<-ctx.Done()
}
