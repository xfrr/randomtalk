package main

import (
	"context"
	"os"
	"os/signal"

	chatcontext "github.com/xfrr/randomtalk/internal/chat"
	"github.com/xfrr/randomtalk/internal/shared/env"
)

var (
	// ServiceVersion is configurable using go build -ldflags "-X main.ServiceVersion=..."
	// or by setting the SERVICE_VERSION environment variable.
	ServiceVersion = env.GetWithDefault("SERVICE_VERSION", "development")
)

func main() {
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelSignal()

	// initialize chat service
	svc := chatcontext.MustInitService(chatcontext.ServiceVersion(ServiceVersion))
	svc.Start(ctx)
}
