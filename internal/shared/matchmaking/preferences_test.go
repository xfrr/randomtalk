package matchmaking_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

type fakeUser struct {
	age   int32
	g     gender.Gender
	prefs matchmaking.Preferences
}

func (u fakeUser) ID() string            { return "" }
func (u fakeUser) Age() int32            { return u.age }
func (u fakeUser) Gender() gender.Gender { return u.g }
func (u fakeUser) Preferences() matchmaking.Preferences {
	return u.prefs
}

func TestDefaultPreferences(t *testing.T) {
	p := matchmaking.DefaultPreferences()
	assert.Equal(t, matchmaking.MinAllowedAge, p.MinAge)
	assert.Equal(t, matchmaking.MaxAllowedAge, p.MaxAge)
	assert.True(t, p.Gender.IsUnspecified())
	assert.Nil(t, p.Interests)
}

func TestWithMinAndMaxAge(t *testing.T) {
	p := matchmaking.DefaultPreferences()

	// Below minimum
	p1 := p.WithMinAge(10)
	assert.Equal(t, matchmaking.MinAllowedAge, p1.MinAge)

	// Above minimum
	p2 := p.WithMinAge(25)
	assert.Equal(t, int32(25), p2.MinAge)

	// Zero max → clamp to max
	p3 := p.WithMaxAge(0)
	assert.Equal(t, matchmaking.MaxAllowedAge, p3.MaxAge)

	// Too large max → clamp to max
	p4 := p.WithMaxAge(200)
	assert.Equal(t, matchmaking.MaxAllowedAge, p4.MaxAge)

	// Valid max
	p5 := p.WithMaxAge(50)
	assert.Equal(t, int32(50), p5.MaxAge)
}

func TestWithGender(t *testing.T) {
	p := matchmaking.DefaultPreferences()
	p1 := p.WithGender(gender.Male)
	assert.True(t, p1.Gender.IsMale())

	// Unspecified leaves unchanged
	p2 := p1.WithGender(gender.Unspecified)
	assert.True(t, p2.Gender.IsMale())
}

func TestWithInterests(t *testing.T) {
	p := matchmaking.DefaultPreferences()
	p1 := p.WithInterests([]string{"x", "y"})
	assert.Equal(t, []string{"x", "y"}, p1.Interests)

	// empty slice no-op
	p2 := p1.WithInterests([]string{})
	assert.Equal(t, p1.Interests, p2.Interests)
}

func TestJSONRoundTripAndDefaults(t *testing.T) {
	// Marshal omits zero-value fields
	p := matchmaking.DefaultPreferences().
		WithMinAge(20).
		WithGender(gender.Female).
		WithInterests([]string{"a"})
	b, err := json.Marshal(p)
	require.NoError(t, err)
	var got matchmaking.Preferences
	err = json.Unmarshal(b, &got)
	require.NoError(t, err)

	// Missing JSON fields default correctly
	assert.Equal(t, int32(20), got.MinAge)
	assert.Equal(t, matchmaking.MaxAllowedAge, got.MaxAge)
	assert.True(t, got.Gender.IsFemale())
	assert.Equal(t, []string{"a"}, got.Interests)

	// Invalid JSON yields ErrInvalidPreferences
	var bad matchmaking.Preferences
	err = json.Unmarshal([]byte(`{"min_age":"x"}`), &bad)
	assert.ErrorIs(t, err, matchmaking.ErrInvalidPreferences)
}

func TestIsSatisfiedBy(t *testing.T) {
	basePrefs := matchmaking.DefaultPreferences().
		WithMinAge(18).
		WithMaxAge(30).
		WithGender(gender.Male).
		WithInterests([]string{"jazz", "rock"})

	// user meets all
	u1 := fakeUser{
		age:   25,
		g:     gender.Male,
		prefs: matchmaking.DefaultPreferences().WithInterests([]string{"rock", "pop"}),
	}
	assert.True(t, basePrefs.IsSatisfiedBy(u1))

	// age too low
	u2 := u1
	u2.age = 17
	assert.False(t, basePrefs.IsSatisfiedBy(u2))

	// wrong gender
	u3 := u1
	u3.g = gender.Female
	assert.False(t, basePrefs.IsSatisfiedBy(u3))

	// no overlapping interest
	u4 := u1
	u4.prefs = matchmaking.DefaultPreferences().WithInterests([]string{"classical"})
	assert.False(t, basePrefs.IsSatisfiedBy(u4))

	// user has no interests → fails
	u5 := u1
	u5.prefs = matchmaking.DefaultPreferences()
	assert.False(t, basePrefs.IsSatisfiedBy(u5))
}
