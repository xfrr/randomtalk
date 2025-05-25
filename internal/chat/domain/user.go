package chatdomain

import (
	domainerr "github.com/xfrr/randomtalk/internal/shared/domain"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	geo "github.com/xfrr/randomtalk/internal/shared/location"
)

const (
	// MinUserAge is the minimum age allowed for a User.
	MinUserAge = 18
	// MaxUserAge is the maximum age allowed for a User.
	MaxUserAge = 120
)

var (
	// ErrUserAgeTooLow is returned when the User age is too low.
	ErrUserAgeTooLow = domainerr.New("user age is too low")
	// ErrUserAgeTooHigh is returned when the User age is too high.
	ErrUserAgeTooHigh = domainerr.New("user age is too high")
)

type NewUserOption func(u *User)

// WithLocation sets the current location of the User.
func WithLocation(location *geo.Location) NewUserOption {
	return func(u *User) {
		u.location = location
	}
}

// User represents a user in the Chat bounded context.
type User struct {
	id               ID
	nickname         string
	age              int32
	location         *geo.Location
	gender           gender.Gender
	matchPreferences MatchPreferences
}

// ID returns the User ID.
func (u User) ID() ID {
	return u.id
}

// Nickname returns the User nickname.
func (u User) Nickname() string {
	return u.nickname
}

// Age returns the User age.
func (u User) Age() int32 {
	return u.age
}

// Location returns the User location.
func (u User) Location() *geo.Location {
	return u.location
}

// Gender returns the gender of the User.
func (u User) Gender() gender.Gender {
	return u.gender
}

// MatchPreferences returns the User match preferences.
func (u User) MatchPreferences() MatchPreferences {
	return u.matchPreferences
}

func (u User) validate() error {
	if u.age < MinUserAge {
		return ErrUserAgeTooLow
	}

	if u.age > MaxUserAge {
		return ErrUserAgeTooHigh
	}
	return nil
}

// NewUser creates a new User instance.
func NewUser(
	id ID,
	nickname string,
	age int32,
	gender gender.Gender,
	matchPreferences MatchPreferences,
	opts ...NewUserOption,
) (User, error) {
	user := User{
		id:               id,
		nickname:         nickname,
		age:              age,
		gender:           gender,
		matchPreferences: matchPreferences,
	}

	for _, opt := range opts {
		opt(&user)
	}

	if err := user.validate(); err != nil {
		return User{}, err
	}

	return user, nil
}
