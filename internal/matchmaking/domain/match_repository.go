package matchdomain

import (
	"context"

	domainerror "github.com/xfrr/randomtalk/internal/shared/domain-error"
)

var (
	ErrMatchAlreadyExists = domainerror.New("match already exists")
	ErrMatchNotFound      = domainerror.New("match not found")
	ErrNoActiveUsers      = domainerror.New("no active users available")
)

type MatchRepository interface {
	// Save persists a match .
	Save(ctx context.Context, match *Match) error

	// FindByID retrieves a match  by its ID.
	FindByID(ctx context.Context, id string) (*Match, error)

	// FindLastByUserID retrieves the last match for the given user ID.
	FindLastByUserID(ctx context.Context, userID string) (*Match, error)

	// Exists checks if a match  with the given ID exists.
	Exists(ctx context.Context, id string) (bool, error)
}
