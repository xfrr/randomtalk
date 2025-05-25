package matchmaking

import "github.com/xfrr/randomtalk/internal/shared/gender"

// User represents an entity that can be matched.
type User interface {
	ID() string
	Age() int32
	Gender() gender.Gender
	Preferences() Preferences
}
