// package gender provides a strongly-typed Gender enumeration with
// robust string, text, and JSON marshalling/unmarshalling.
package gender

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidGender is returned when unmarshalling an empty or invalid gender string.
	ErrInvalidGender = errors.New("invalid gender")
)

// Gender is an enumeration of gender values.
type Gender int

const (
	Unspecified Gender = iota
	Female
	Male
)

var (
	// genderNames holds the canonical string values for each Gender.
	genderNames = []string{
		"unspecified",
		"female",
		"male",
	}
	// genderMap maps lowercase string values to their Gender.
	genderMap = func() map[string]Gender {
		m := make(map[string]Gender, len(genderNames))
		for i, name := range genderNames {
			m[name] = Gender(i)
		}
		return m
	}()
)

// String returns the canonical lowercase string for a Gender.
// If g is out of range, it returns "unspecified".
func (g Gender) String() string {
	idx := int(g)
	if idx < 0 || idx >= len(genderNames) {
		return genderNames[0]
	}
	return genderNames[idx]
}

// Parse returns the Gender corresponding to s (case-insensitive).
// Unknown values yield GenderUnspecified.
func Parse(s string) Gender {
	if g, ok := genderMap[strings.ToLower(s)]; ok {
		return g
	}
	return Unspecified
}

// IsValid reports whether g is one of the defined constants.
func (g Gender) IsValid() bool {
	_, ok := genderMap[g.String()]
	return ok
}

// Convenience checks.
func (g Gender) IsMale() bool        { return g == Male }
func (g Gender) IsFemale() bool      { return g == Female }
func (g Gender) IsUnspecified() bool { return g == Unspecified }
func (g Gender) Is(v Gender) bool    { return g == v }

// MarshalText implements encoding.TextMarshaler.
func (g Gender) MarshalText() ([]byte, error) {
	return []byte(g.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (g *Gender) UnmarshalText(text []byte) error {
	s := strings.ToLower(string(text))
	if s == "" {
		return ErrInvalidGender
	}
	if val, ok := genderMap[s]; ok {
		*g = val
		return nil
	}
	*g = Unspecified
	return nil
}

// MarshalJSON implements json.Marshaler via text marshaling.
func (g Gender) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// UnmarshalJSON implements json.Unmarshaler via text unmarshaling.
func (g *Gender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("gender must be a string, got %s", string(data))
	}
	return g.UnmarshalText([]byte(s))
}
