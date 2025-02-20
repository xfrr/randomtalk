package matchdomain

import "context"

// MatchmakingProcessor defines how users request matches and
// how a background worker processes them.
type MatchmakingProcessor interface {
	// ProcessMatchRequest stores and enqueues the user.
	ProcessMatchRequest(ctx context.Context, user User) error
}

// StableMatchFinder defines the interface for a stable matching algorithm.
type StableMatchFinder interface {
	// FindStableMatches runs the stable matching algorithm on two sets of users.
	// The algorithm should return a list of indexes from setA to setB, where
	// matches[a] = b means that user setA[a] is matched with user setB[b].
	//
	// The algorithm should return nil if no matches are possible.
	FindStableMatches(setA, setB []*User) []int
}
