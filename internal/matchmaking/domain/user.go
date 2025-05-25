package matchdomain

import (
	"encoding/json"
	"strings"

	domain_error "github.com/xfrr/randomtalk/internal/shared/domain"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

// Domain-level errors
var (
	ErrMatchAlreadyMatchedCandidate = domain_error.New("match already matched candidate")
	ErrRequesterAndCandidateAreSame = domain_error.New("requester and candidate are the same")
	ErrCandidateDoesNotPreferences  = domain_error.New("candidate does not match preferences")
	ErrUserNotWaitingForMatch       = domain_error.New("user is not waiting for a match")
	ErrUserNotFound                 = domain_error.New("user not found")
)

// UserStatus represents the status of a user in the matchmaking process.
type UserStatus int

const (
	Waiting UserStatus = iota
	Matched
	Rejected
)

func (s UserStatus) String() string {
	switch s {
	case Waiting:
		return "waiting"
	case Matched:
		return "matched"
	case Rejected:
		return "rejected"
	}
	return "waiting"
}

// ParseUserStatus returns the UserStatus for a given string, defaulting to Waiting.
func ParseUserStatus(str string) UserStatus {
	switch strings.ToLower(str) {
	case "waiting":
		return Waiting
	case "matched":
		return Matched
	case "rejected":
		return Rejected
	default:
		return Waiting
	}
}

// MarshalText implements encoding.TextMarshaler for UserStatus.
func (s UserStatus) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for UserStatus.
func (s *UserStatus) UnmarshalText(text []byte) error {
	*s = ParseUserStatus(string(text))
	return nil
}

// User is a domain entity representing a user who can be matched.
type User struct {
	id     string
	age    int32
	gender gender.Gender
	prefs  matchmaking.Preferences
	status UserStatus
}

// NewUser constructs a new User with default status=Waiting.
func NewUser(
	id string,
	age int32,
	g gender.Gender,
	preferences matchmaking.Preferences,
) *User {
	return &User{
		id:     id,
		age:    age,
		gender: g,
		prefs:  preferences,
		status: Waiting,
	}
}

// ID returns the user's unique identifier.
func (u User) ID() string { return u.id }

// Age returns the user's age.
func (u User) Age() int32 { return u.age }

// Gender returns the user's gender.
func (u User) Gender() gender.Gender { return u.gender }

// Preferences returns the user's match preferences.
func (u User) Preferences() matchmaking.Preferences { return u.prefs }

// Status returns the user's current status.
func (u User) Status() UserStatus { return u.status }

// Waiting reports whether the user is waiting for a match.
func (u User) Waiting() bool { return u.status == Waiting }

// Matched reports whether the user has been matched.
func (u User) Matched() bool { return u.status == Matched }

// Rejected reports whether the user has been rejected.
func (u User) Rejected() bool { return u.status == Rejected }

// SetStatus transitions the user to a new status.
func (u *User) SetStatus(status UserStatus) { u.status = status }

// MarshalJSON serializes the User to JSON, preserving encapsulation.
func (u User) MarshalJSON() ([]byte, error) {
	type dto struct {
		ID          string                  `json:"id"`
		Age         int32                   `json:"age"`
		Gender      gender.Gender           `json:"gender"`
		Preferences matchmaking.Preferences `json:"preferences"`
		Status      UserStatus              `json:"status"`
	}
	return json.Marshal(dto{
		ID:          u.id,
		Age:         u.age,
		Gender:      u.gender,
		Preferences: u.prefs,
		Status:      u.status,
	})
}

// UnmarshalJSON deserializes JSON into a User, preserving encapsulation.
func (u *User) UnmarshalJSON(data []byte) error {
	type dto struct {
		ID          string                  `json:"id"`
		Age         int32                   `json:"age"`
		Gender      gender.Gender           `json:"gender"`
		Preferences matchmaking.Preferences `json:"preferences"`
		Status      UserStatus              `json:"status"`
	}
	var d dto
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	u.id = d.ID
	u.age = d.Age
	u.gender = d.Gender
	u.prefs = d.Preferences
	u.status = d.Status
	return nil
}
