package matchgrpc

import (
	"context"
	"net"

	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/location"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"

	matchcommands "github.com/xfrr/randomtalk/internal/matchmaking/application/commands"
	matchqueries "github.com/xfrr/randomtalk/internal/matchmaking/application/queries"
	matchpb "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/matchmaking/v1"

	"github.com/xfrr/go-cqrsify/cqrs"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	pbreflection "google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var _ matchpb.MatchMakingServiceServer = &Server{}

type Server struct {
	*grpc.Server
	matchpb.UnimplementedMatchMakingServiceServer

	cmdbus   matchcommands.CommandBus
	querybus matchqueries.QueryBus
}

func (s *Server) FindMatch(
	ctx context.Context,
	in *matchpb.FindMatchRequest,
) (*matchpb.FindMatchResponse, error) {
	cmd := matchcommands.MatchUserWithPreferencesCommand{
		UserID:     in.GetUserId(),
		UserAge:    in.GetUserAge(),
		UserGender: toGender(in.GetUserGender()),
		UserPreferences: matchmaking.DefaultPreferences().
			WithMinAge(in.GetMatchPreferences().GetMinAge()).
			WithMaxAge(in.GetMatchPreferences().GetMaxAge()).
			WithGender(toGender(in.GetMatchPreferences().GetGender())).
			WithInterests(in.GetMatchPreferences().GetInterests()),
	}

	res, err := cqrs.Dispatch[*matchcommands.MatchUserWithPreferencesResponse](ctx, s.cmdbus, cmd)
	if err != nil {
		return nil, handleError(err)
	}

	return &matchpb.FindMatchResponse{
		MatchId: res.MatchID,
	}, nil
}

func handleError(err error) error {
	return status.Error(status.Code(err), err.Error())
}

func NewServer(
	matchCommandBus matchcommands.CommandBus,
	matchQueryBus matchqueries.QueryBus,
) (*Server, func(), error) {
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	s := &Server{
		Server:   grpcServer,
		cmdbus:   matchCommandBus,
		querybus: matchQueryBus,
	}

	matchpb.RegisterMatchMakingServiceServer(grpcServer, s)

	// register reflection service on gRPC server
	pbreflection.Register(grpcServer)

	return s, grpcServer.GracefulStop, nil
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

func toGender(g matchpb.Gender) gender.Gender {
	switch g {
	case matchpb.Gender_GENDER_UNSPECIFIED:
		return gender.Unspecified
	case matchpb.Gender_GENDER_MALE:
		return gender.Male
	case matchpb.Gender_GENDER_FEMALE:
		return gender.Female
	default:
		return gender.Unspecified
	}
}

func toLocation(latlng *matchpb.LatLng) *location.Location {
	if latlng == nil {
		return nil
	}

	userLocation := location.New(
		latlng.GetLatitude(),
		latlng.GetLongitude(),
	)
	return &userLocation
}
