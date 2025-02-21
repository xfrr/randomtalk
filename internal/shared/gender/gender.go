package gender

import (
	"errors"
	"strings"
)

var (
	ErrUnmarshalJSON = errors.New("failed to unmarshal gender from json")
)

type Gender int

const (
	GenderUnspecified Gender = iota
	GenderFemale
	GenderMale
)

var GenderStrings = map[string]Gender{
	"unspecified": GenderUnspecified,
	"female":      GenderFemale,
	"male":        GenderMale,
}

var GenderValues = map[Gender]string{
	GenderUnspecified: "unspecified",
	GenderFemale:      "female",
	GenderMale:        "male",
}

func (g Gender) String() string {
	return GenderValues[g]
}

func (g Gender) IsValid() bool {
	_, ok := GenderValues[g]
	return ok
}

func (g Gender) IsMale() bool {
	return g == GenderMale
}

func (g Gender) IsFemale() bool {
	return g == GenderFemale
}

func (g Gender) IsUnspecified() bool {
	return g == GenderUnspecified
}

func (g Gender) Is(v Gender) bool {
	return g == v
}

func (g *Gender) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	if str == "" {
		return ErrUnmarshalJSON
	}

	switch str {
	case "male", "gender_male":
		*g = GenderMale
	case "female", "gender_female":
		*g = GenderFemale
	default:
		*g = GenderUnspecified
	}
	return nil
}

func (g *Gender) MarshalJSON() ([]byte, error) {
	return []byte("\"" + g.String() + "\""), nil
}

func ParseString(s string) Gender {
	return GenderStrings[strings.ToLower(s)]
}
