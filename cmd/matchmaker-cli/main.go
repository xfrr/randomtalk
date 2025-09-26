package main

import (
	"context"

	"github.com/spf13/cobra"

	matchsessioncli "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/cli"
)

var RootCmd = &cobra.Command{
	Use:   "randomtalk",
	Short: "Random chat CLI",
}

var (
	grpcAddr = RootCmd.Flags().String("grpc-addr", "localhost:50000", "gRPC server address")
)

func main() {
	ctx := context.Background()

	// add cobra commands
	RootCmd.AddCommand(matchsessioncli.NewMatchSessionCobraCommand())

	err := RootCmd.ExecuteContext(ctx)
	if err != nil {
		panic(err)
	}
}
