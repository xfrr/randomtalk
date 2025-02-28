package matchmakingtrace

import (
	"context"
	"time"

	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ matchdomain.MatchRepository = (*TraceableMatchRepository)(nil)

// TraceableMatchRepository is a traceable implementation of the MatchRepository interface.
type TraceableMatchRepository struct {
	repo   matchdomain.MatchRepository
	tracer trace.Tracer
}

// WrapMatchRepository creates a new MatchRepositoryTraceable.
func WrapMatchRepository(
	repo matchdomain.MatchRepository,
	traceProvider trace.TracerProvider,
) *TraceableMatchRepository {
	return &TraceableMatchRepository{
		repo:   repo,
		tracer: traceProvider.Tracer("randomtalk_match_repository"),
	}
}

// Save creates a new match.
func (r *TraceableMatchRepository) Save(ctx context.Context, match *matchdomain.Match) error {
	ctx, span := r.tracer.Start(
		ctx,
		"Save",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(
			attribute.String("match_id", match.ID()),
		),
	)
	defer span.End()

	err := r.repo.Save(ctx, match)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// FindByID finds a match by its ID.
func (r *TraceableMatchRepository) FindByID(ctx context.Context, matchID string) (*matchdomain.Match, error) {
	ctx, span := r.tracer.Start(
		ctx,
		"FindByID",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(
			attribute.String("match_id", matchID),
		),
	)
	defer span.End()

	match, err := r.repo.FindByID(ctx, matchID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return match, err
}

// Exists checks if a match exists.
func (r *TraceableMatchRepository) Exists(ctx context.Context, matchID string) (bool, error) {
	ctx, span := r.tracer.Start(
		ctx,
		"Exists",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(
			attribute.String("match_id", matchID),
		),
	)
	defer span.End()

	exists, err := r.repo.Exists(ctx, matchID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return exists, err
}
