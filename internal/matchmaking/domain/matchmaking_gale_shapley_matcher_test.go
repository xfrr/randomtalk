package matchdomain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

func TestGaleShapleyStableMatcher_FindStableMatches(t *testing.T) {
	t.Run("store returns error", func(t *testing.T) {
		matcher := domain.NewGaleShapleyStableMatcher()

		matches := matcher.FindStableMatches(nil, nil)
		assert.Nil(t, matches)
	})

	t.Run("unrestricted preferences scenario", func(t *testing.T) {
		// Both B1 and B2 have default preferences
		userB1 := domain.NewUser("B1", 25, gender.GenderUnspecified, matchmaking.DefaultPreferences())
		userB2 := domain.NewUser("B2", 30, gender.GenderUnspecified, matchmaking.DefaultPreferences())

		matcher := domain.NewGaleShapleyStableMatcher()
		userA1 := domain.NewUser("A1", 20, gender.GenderUnspecified, matchmaking.DefaultPreferences())

		matches := matcher.FindStableMatches([]*domain.User{&userA1}, []*domain.User{&userB1, &userB2})
		require.Len(t, matches, 1, "expecting a match for A1")
		assert.Equal(t, 0, matches[0], "A1 should be matched with B1")
	})

	t.Run("multiple A and B with partial compatibility", func(t *testing.T) {
		// B1 only compatible with A1, B2 only with A2
		prefsB1 := matchmaking.DefaultPreferences()
		prefsB2 := matchmaking.DefaultPreferences()
		userB1 := domain.NewUser("B1", 25, gender.GenderUnspecified, prefsB1)
		userB2 := domain.NewUser("B2", 22, gender.GenderUnspecified, prefsB2)

		matcher := domain.NewGaleShapleyStableMatcher()

		// A1 and A2 have default preferences, assuming no additional constraints
		userA1 := domain.NewUser("A1", 20, gender.GenderUnspecified, matchmaking.DefaultPreferences())
		userA2 := domain.NewUser("A2", 24, gender.GenderUnspecified, matchmaking.DefaultPreferences())

		matches := matcher.FindStableMatches([]*domain.User{&userA1, &userA2}, []*domain.User{&userB1, &userB2})
		require.Len(t, matches, 2, "expecting matches for A1 and A2")
		assert.Equal(t, 0, matches[0], "A1 should be matched with B1")
		assert.Equal(t, 1, matches[1], "A2 should be matched with B2")
	})

	t.Run("no compatibility found", func(t *testing.T) {
		// B1 and B2
		userB1 := domain.NewUser("B1", 25, gender.GenderUnspecified, matchmaking.DefaultPreferences())
		userB2 := domain.NewUser("B2", 60, gender.GenderMale, matchmaking.DefaultPreferences())

		matcher := domain.NewGaleShapleyStableMatcher()

		// A1 and A2 have default preferences, assuming no additional constraints
		userA1 := domain.NewUser("A1", 20, gender.GenderMale, matchmaking.DefaultPreferences().WithMaxAge(18))
		userA2 := domain.NewUser("A2", 24, gender.GenderFemale, matchmaking.DefaultPreferences().WithMaxAge(18))

		matches := matcher.FindStableMatches([]*domain.User{&userA1, &userA2}, []*domain.User{&userB1, &userB2})
		require.Len(t, matches, 2, "expecting matches for A1 and A2")
		assert.Equal(t, -1, matches[0], "A1 should not be matched")
		assert.Equal(t, -1, matches[1], "A2 should not be matched")
	})
}
