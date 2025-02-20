package matchmakingtrace

import (
	"context"
	"time"

	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ matchdomain.MatchmakingProcessor = (*TraceableMatchmakingService)(nil)

// TraceableMatchmakingService is a traceable implementation of the MatchmakingService interface.
type TraceableMatchmakingService struct {
	service matchdomain.MatchmakingProcessor
	tracer  trace.Tracer
}

// WrapMatchmakingService creates a new MatchmakingServiceTraceable.
func WrapMatchmakingService(
	service matchdomain.MatchmakingProcessor,
	tracer trace.TracerProvider,
) *TraceableMatchmakingService {
	return &TraceableMatchmakingService{
		service: service,
		tracer:  tracer.Tracer("randomtalk_matchmaking_service"),
	}
}

// ProcessMatchRequest matches a user with other users based on their preferences.
func (s *TraceableMatchmakingService) ProcessMatchRequest(ctx context.Context, user matchdomain.User) error {
	ctx, span := s.tracer.Start(
		ctx, "RequestMatch",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(
			attribute.String("user_id", user.ID()),
			attribute.Int("user_age", user.Age()),
			attribute.String("user_gender", user.Gender().String()),
		))
	defer span.End()

	err := s.service.ProcessMatchRequest(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
