package matchmaking

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/xfrr/randomtalk/internal/shared/gender"
)

// ErrInvalidPreferences is returned when JSON unmarshalling fails.
var ErrInvalidPreferences = errors.New("invalid match preferences")

const (
	MinAllowedAge int32 = 18
	MaxAllowedAge int32 = 99
)

// Preferences holds criteria for matching users.
type Preferences struct {
	MinAge    int32         `json:"min_age"`
	MaxAge    int32         `json:"max_age"`
	Gender    gender.Gender `json:"gender,omitempty"`
	Interests []string      `json:"interests,omitempty"`
}

// DefaultPreferences returns a Preferences with sane defaults.
func DefaultPreferences() Preferences {
	return Preferences{
		MinAge: MinAllowedAge,
		MaxAge: MaxAllowedAge,
		// Gender zero value is Unspecified → omitted in JSON
		// Interests zero value is nil → omitted in JSON
	}
}

// WithMinAge returns a copy with MinAge clamped to [MinAllowedAge, ∞).
func (p Preferences) WithMinAge(min int32) Preferences {
	if min < MinAllowedAge {
		min = MinAllowedAge
	}
	p.MinAge = min
	return p
}

// WithMaxAge returns a copy with MaxAge clamped to (0, MaxAllowedAge].
func (p Preferences) WithMaxAge(max int32) Preferences {
	if max <= 0 || max > MaxAllowedAge {
		max = MaxAllowedAge
	}
	p.MaxAge = max
	return p
}

// WithGender returns a copy with Gender set (ignores Unspecified).
func (p Preferences) WithGender(g gender.Gender) Preferences {
	if !g.IsUnspecified() {
		p.Gender = g
	}
	return p
}

// WithInterests returns a copy with a non-empty interests slice.
func (p Preferences) WithInterests(interests []string) Preferences {
	if len(interests) == 0 {
		return p
	}
	cp := make([]string, len(interests))
	copy(cp, interests)
	p.Interests = cp
	return p
}

// UnmarshalJSON applies defaults when fields are missing or zero.
func (p *Preferences) UnmarshalJSON(data []byte) error {
	type alias Preferences
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return ErrInvalidPreferences
	}
	// defaults
	if tmp.MinAge == 0 {
		tmp.MinAge = MinAllowedAge
	}
	if tmp.MaxAge == 0 {
		tmp.MaxAge = MaxAllowedAge
	}
	*p = Preferences(tmp)
	return nil
}

// String implements fmt.Stringer for debug output.
func (p Preferences) String() string {
	parts := []string{
		fmt.Sprintf("MinAge: %d", p.MinAge),
		fmt.Sprintf("MaxAge: %d", p.MaxAge),
	}
	if !p.Gender.IsUnspecified() {
		parts = append(parts, "Gender: "+p.Gender.String())
	}
	if len(p.Interests) > 0 {
		parts = append(parts, "Interests: ["+strings.Join(p.Interests, ", ")+"]")
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// IsSatisfiedBy reports whether a User meets all criteria.
func (p Preferences) IsSatisfiedBy(u User) bool {
	age := u.Age()
	if age < p.MinAge || age > p.MaxAge {
		return false
	}
	if !p.Gender.IsUnspecified() && !p.Gender.Is(u.Gender()) {
		return false
	}
	if len(p.Interests) > 0 {
		userInterests := u.Preferences().Interests
		if len(userInterests) == 0 {
			return false
		}
		match := false
		for _, want := range p.Interests {
			for _, have := range userInterests {
				if want == have {
					match = true
					break
				}
			}
			if match {
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}
