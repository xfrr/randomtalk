package matchmakinginmemory

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
)

var _ matchdomain.UserStore = (*UserStore)(nil)

// UserStore implements matchdomain.UserStore using an in-memory concurrent implementation.
type UserStore struct {
	usersIndex sync.Map
	logger     *zerolog.Logger
}

// NewUserStore initializes an in-memory user store.
func NewUserStore(logger *zerolog.Logger) *UserStore {
	return &UserStore{
		logger: logger,
	}
}

// AddUser adds a user to the in-memory store.
func (us *UserStore) AddUser(_ context.Context, user matchdomain.User) error {
	us.usersIndex.Store(user.ID(), user)
	return nil
}

// FindByID finds a user by ID in the in-memory store.
func (us *UserStore) FindByID(_ context.Context, userID string) (*matchdomain.User, error) {
	if user, ok := us.usersIndex.Load(userID); ok {
		u, _ := user.(matchdomain.User)
		return &u, nil
	}
	return nil, matchdomain.ErrUserNotFound
}

// GetAll retrieves all users from the in-memory store.
func (us *UserStore) GetAll(_ context.Context) ([]*matchdomain.User, error) {
	var users []*matchdomain.User
	us.usersIndex.Range(func(_, value interface{}) bool {
		u, _ := value.(matchdomain.User)
		users = append(users, &u)
		return true
	})
	return users, nil
}

// FindCompatibleUser finds a compatible user in the in-memory store.
func (us *UserStore) FindUserByPreferences(_ context.Context, user *matchdomain.User) (*matchdomain.User, error) {
	var compatibleUser matchdomain.User
	us.usersIndex.Range(func(_, value interface{}) bool {
		u, _ := value.(matchdomain.User)
		if user.ID() != u.ID() && user.MatchPreferences().IsSatisfiedBy(&u) {
			compatibleUser = u
			return false
		}

		return true
	})

	if compatibleUser.ID() == "" {
		return nil, matchdomain.ErrUserNotFound
	}

	return &compatibleUser, nil
}

// RemoveUsers removes users from the in-memory store.
func (us *UserStore) RemoveUsers(ctx context.Context, userIDs ...string) error {
	// check if users exist
	for _, userID := range userIDs {
		if _, err := us.FindByID(ctx, userID); err != nil {
			return err
		}
	}

	// remove users
	for _, userID := range userIDs {
		us.usersIndex.Delete(userID)
	}

	return nil
}
