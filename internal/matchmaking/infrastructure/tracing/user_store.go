package matchmakingtrace

import (
	"context"

	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var _ matchdomain.UserStore = (*TraceableUserStore)(nil)

// TraceableUserStore is a traceable implementation of the WaitingUserStore interface.
type TraceableUserStore struct {
	store  matchdomain.UserStore
	tracer trace.Tracer
}

// WrapUserStore creates a new UserStoreTraceable.
func WrapUserStore(
	userStore matchdomain.UserStore,
	traceProvider trace.TracerProvider,
) *TraceableUserStore {
	return &TraceableUserStore{
		store:  userStore,
		tracer: traceProvider.Tracer("randomtalk_matchmaking_user_store"),
	}
}

// AddUser adds a user to the matchmaking queue.
func (us *TraceableUserStore) AddUser(ctx context.Context, user matchdomain.User) error {
	ctx, span := us.tracer.Start(ctx, "randomtalk/matchmaking/user_store/AddUser")
	defer span.End()

	err := us.store.AddUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to add user")
	}

	return err
}

// GetAll retrieves all users from the matchmaking queue.
func (us *TraceableUserStore) GetAll(ctx context.Context) ([]*matchdomain.User, error) {
	ctx, span := us.tracer.Start(ctx, "randomtalk/matchmaking/user_store/GetAll")
	defer span.End()

	users, err := us.store.GetAll(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get all users")
	}

	return users, err
}

func (us *TraceableUserStore) RemoveUsers(ctx context.Context, userIDs ...string) error {
	ctx, span := us.tracer.Start(ctx, "randomtalk/matchmaking/user_store/RemoveUsers")
	defer span.End()

	err := us.store.RemoveUsers(ctx, userIDs...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to remove users")
	}

	return err
}
