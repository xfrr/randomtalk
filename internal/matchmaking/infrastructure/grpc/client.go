package matchgrpc

import (
	"context"

	matchpb "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/matchmaking/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ matchpb.MatchMakingServiceClient = &MatchManagerClient{}

type MatchManagerClient struct {
	conn matchpb.MatchMakingServiceClient
}

func (c *MatchManagerClient) GetMatch(ctx context.Context, in *matchpb.GetMatchRequest, opts ...grpc.CallOption) (*matchpb.GetMatchResponse, error) {
	res, err := c.conn.GetMatch(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *MatchManagerClient) FindMatch(ctx context.Context, in *matchpb.FindMatchRequest, opts ...grpc.CallOption) (*matchpb.FindMatchResponse, error) {
	res, err := c.conn.FindMatch(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// NewClient finds a new MatchManagerClient.
func NewClient(addr string) (*MatchManagerClient, func(), error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	client := matchpb.NewMatchMakingServiceClient(conn)
	return &MatchManagerClient{conn: client}, func() { conn.Close() }, nil
}
