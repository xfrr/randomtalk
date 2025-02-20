package matchdomain

import "context"

// UserStore defines the behavior of a user store.
type UserStore interface {
	// AddUser adds a user to the store.
	AddUser(ctx context.Context, user User) error

	// GetAll returns all users in the store.
	GetAll(ctx context.Context) ([]*User, error)

	// RemoveUsers removes a user from the store.
	RemoveUsers(ctx context.Context, userID ...string) error
}
