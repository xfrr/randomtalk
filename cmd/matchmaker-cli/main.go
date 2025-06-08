package main

import (
	"context"

	"github.com/spf13/cobra"

	matchsessioncli "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/cli"
	matchsessiongrpc "github.com/xfrr/randomtalk/internal/matchmaking/infrastructure/grpc"
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

	grpcClient, grpcClose, err := matchsessiongrpc.NewClient(*grpcAddr)
	if err != nil {
		panic(err)
	}
	defer grpcClose()

	// add cobra commands
	RootCmd.AddCommand(matchsessioncli.NewMatchSessionCobraCommand(grpcClient))

	err = RootCmd.ExecuteContext(ctx)
	if err != nil {
		panic(err)
	}
}
