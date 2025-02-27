package matchmaking

import (
	"github.com/xfrr/randomtalk/internal/shared/gender"
)

// User represents a user in the matchmaking system.
type User interface {
	// ID returns the user's unique identifier.
	ID() string

	// Age returns the user's age.
	Age() int

	// Gender returns the user's gender
	Gender() gender.Gender

	// MatchPreferences returns the user's match preferences.
	MatchPreferences() MatchPreferences
}
