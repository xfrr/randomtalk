package chatdomain

import "context"

// MatchRequester defines the interface for requesting a match+ for a given ChatSession.
type MatchRequester interface {
	RequestMatch(ctx context.Context, cs *ChatSession) error
}
