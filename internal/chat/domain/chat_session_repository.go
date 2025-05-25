package chatdomain

import (
	"context"

	domainerror "github.com/xfrr/randomtalk/internal/shared/domain"
)

var (
	ErrChatSessionNotFound      = domainerror.New("chat session not found")
	ErrChatSessionAlreadyExists = domainerror.New("chat session already exists with the given ID")
)

type ChatSessionRepository interface {
	// Save persists the state of the ChatSession Aggregate.
	Save(ctx context.Context, cs *ChatSession) error

	// FindByID retrieves a ChatSession by its unique identifier.
	FindByID(ctx context.Context, id string) (*ChatSession, error)

	// Exists checks if a ChatSession with the given ID exists.
	Exists(ctx context.Context, id string) (bool, error)
}
