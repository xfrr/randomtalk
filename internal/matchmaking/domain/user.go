package matchdomain

import (
	"encoding/json"

	domain_error "github.com/xfrr/randomtalk/internal/shared/domain-error"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

// Domain-level errors
var (
	ErrMatchAlreadyMatchedCandidate     = domain_error.New("match already matched candidate")
	ErrRequesterAndCandidateAreSame     = domain_error.New("requester and candidate are the same")
	ErrCandidateDoesNotMatchPreferences = domain_error.New("candidate does not match preferences")
	ErrUserNotWaitingForMatch           = domain_error.New("user is not waiting for a match")
	ErrUserNotFound                     = domain_error.New("user not found")
)

// UserStatus represents the status of a user in the matchmaking process.
type UserStatus int

func (s UserStatus) String() string {
	return [...]string{"waiting", "matched", "rejected"}[s]
}

func ParseUserStatus(s string) UserStatus {
	switch s {
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

const (
	// Waiting status indicates the user is waiting for a match.
	Waiting UserStatus = iota

	// Matched status indicates the user has been matched.
	Matched

	// Rejected status indicates the user has been rejected.
	Rejected
)

// User is a domain entity representing a user who can be matched.
type User struct {
	id          string
	age         int
	gender      gender.Gender
	preferences *matchmaking.MatchPreferences
	status      UserStatus
}

// NewUser constructs a new user with default version=1 and status=waiting.
func NewUser(
	id string,
	age int,
	g gender.Gender,
	preferences *matchmaking.MatchPreferences,
) User {
	return User{
		id:          id,
		age:         age,
		gender:      g,
		preferences: preferences,
		status:      Waiting,
	}
}

// --- Getters & Setters ---

func (u *User) ID() string {
	return u.id
}

func (u *User) Age() int {
	return u.age
}

func (u *User) Gender() gender.Gender {
	return u.gender
}

func (u User) MatchPreferences() matchmaking.MatchPreferences {
	if u.preferences == nil {
		return *matchmaking.DefaultPreferences()
	}

	return *u.preferences
}

func (u User) MarshalJSON() ([]byte, error) {
	return MarshalUser(&u)
}

func (u *User) UnmarshalJSON(data []byte) error {
	usr, err := UnmarshalUser(data)
	if err != nil {
		return err
	}

	u.id = usr.id
	u.age = usr.age
	u.gender = usr.gender
	u.preferences = usr.preferences
	u.status = usr.status
	return nil
}

// --- User Status ---
func (u *User) Status() UserStatus {
	return u.status
}

func (u *User) SetStatus(status UserStatus) {
	u.status = status
}

func (u *User) Waiting() bool {
	return u.status == Waiting
}

func (u *User) Matched() bool {
	return u.status == Matched
}

func (u *User) Rejected() bool {
	return u.status == Rejected
}

// --- JSON Serialization Helpers ---

func UnmarshalUser(data []byte) (*User, error) {
	type userCopy struct {
		ID          string                        `json:"id"`
		Age         int                           `json:"age"`
		Gender      gender.Gender                 `json:"gender"`
		Preferences *matchmaking.MatchPreferences `json:"preferences,omitempty"`
		Status      UserStatus                    `json:"status,omitempty"`
	}
	var u userCopy
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, err
	}

	return &User{
		id:          u.ID,
		age:         u.Age,
		gender:      u.Gender,
		preferences: u.Preferences,
		status:      u.Status,
	}, nil
}

func MarshalUser(u *User) ([]byte, error) {
	type userCopy struct {
		ID          string                        `json:"id"`
		Age         int                           `json:"age"`
		Gender      gender.Gender                 `json:"gender"`
		Preferences *matchmaking.MatchPreferences `json:"preferences"`
		Status      UserStatus                    `json:"status"`
	}
	cpy := userCopy{
		ID:          u.id,
		Age:         u.age,
		Gender:      u.gender,
		Preferences: u.preferences,
		Status:      u.status,
	}
	return json.Marshal(cpy)
}
