package gender_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/randomtalk/internal/shared/gender"
)

func TestStringAndParse(t *testing.T) {
	cases := []struct {
		input    gender.Gender
		expected string
	}{
		{gender.Unspecified, "unspecified"},
		{gender.Female, "female"},
		{gender.Male, "male"},
		{gender.Gender(999), "unspecified"},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, c.input.String())
		assert.Equal(t, c.expected, gender.Parse(c.expected).String())
	}

	assert.Equal(t, gender.Unspecified, gender.Parse("not_a_gender"))
}

func TestIsValid(t *testing.T) {
	assert.True(t, gender.Female.IsValid())
	assert.True(t, gender.Male.IsValid())
	assert.True(t, gender.Unspecified.IsValid())
	assert.True(t, gender.Gender(42).IsValid())
}

func TestConvenienceChecks(t *testing.T) {
	assert.True(t, gender.Male.IsMale())
	assert.False(t, gender.Male.IsFemale())
	assert.False(t, gender.Male.IsUnspecified())

	assert.True(t, gender.Female.IsFemale())
	assert.True(t, gender.Unspecified.IsUnspecified())
	assert.True(t, gender.Female.Is(gender.Female))
}

func TestMarshalUnmarshalText(t *testing.T) {
	var g gender.Gender

	for _, name := range []string{"female", "male", "unspecified"} {
		err := g.UnmarshalText([]byte(name))
		require.NoError(t, err)
		text, err := g.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, name, string(text))
	}

	err := g.UnmarshalText([]byte(""))
	require.ErrorIs(t, err, gender.ErrInvalidGender)

	err = g.UnmarshalText([]byte("unknown"))
	require.NoError(t, err)
	assert.Equal(t, gender.Unspecified, g)
}

func TestJSONRoundTrip(t *testing.T) {
	type wrapper struct {
		G gender.Gender `json:"G"`
	}

	for _, name := range []string{"female", "male", "unspecified"} {
		w1 := wrapper{G: gender.Parse(name)}
		b, err := json.Marshal(w1)
		require.NoError(t, err)
		var w2 wrapper
		err = json.Unmarshal(b, &w2)
		require.NoError(t, err)
		assert.Equal(t, w1.G, w2.G)
	}

	var w wrapper
	err := json.Unmarshal([]byte(`{"G":42}`), &w)
	require.Error(t, err)

	err = json.Unmarshal([]byte(`{"G":""}`), &w)
	assert.ErrorIs(t, err, gender.ErrInvalidGender)
}
