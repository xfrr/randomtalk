package matchmakingtrace

import (
	"context"

	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ matchdomain.NotificationsChannel = (*TraceableMatchNotificationsChannel)(nil)

// TraceableMatchNotificationsChannel is a traceable implementation of the NotificationChannel interface.
type TraceableMatchNotificationsChannel struct {
	channel matchdomain.NotificationsChannel
	tracer  trace.Tracer
}

// WrapMatchNotificationsChannel creates a new MatchNotificationsChannelTraceable.
func WrapMatchNotificationsChannel(
	channel matchdomain.NotificationsChannel,
	traceProvider trace.TracerProvider,
) *TraceableMatchNotificationsChannel {
	return &TraceableMatchNotificationsChannel{
		channel: channel,
		tracer:  traceProvider.Tracer("randomtalk_matchmaking_notifications_channel"),
	}
}

// Notify notifies the user with the given ID about the match.
func (c *TraceableMatchNotificationsChannel) Notify(ctx context.Context, userToNotifyID string, match *matchdomain.Match) error {
	ctx, span := c.tracer.Start(
		ctx,
		"Notify",
		trace.WithAttributes(
			attribute.String("user_id", userToNotifyID),
			attribute.String("component", "match_notifications_channel"),
		),
	)
	defer span.End()

	err := c.channel.Notify(ctx, userToNotifyID, match)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
