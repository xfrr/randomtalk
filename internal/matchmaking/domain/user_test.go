package matchdomain_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

func TestNewUser(t *testing.T) {
	prefs := matchmaking.Preferences{MinAge: 18, MaxAge: 30, Gender: gender.Female}
	u := matchdomain.NewUser("user1", 25, gender.Male, prefs)

	require.Equal(t, "user1", u.ID())
	require.Equal(t, int32(25), u.Age())
	require.Equal(t, gender.Male, u.Gender())
	require.Equal(t, prefs, u.Preferences())
	require.Equal(t, matchdomain.Waiting, u.Status())
	require.True(t, u.Waiting())
	require.False(t, u.Matched())
	require.False(t, u.Rejected())
}

func TestUserStatusTransitions(t *testing.T) {
	u := matchdomain.NewUser("user2", 22, gender.Female, matchmaking.Preferences{})

	require.True(t, u.Waiting())
	u.SetStatus(matchdomain.Matched)
	require.True(t, u.Matched())
	require.False(t, u.Waiting())
	u.SetStatus(matchdomain.Rejected)
	require.True(t, u.Rejected())
}

func TestUserStatusStringAndParse(t *testing.T) {
	require.Equal(t, "waiting", matchdomain.Waiting.String())
	require.Equal(t, "matched", matchdomain.Matched.String())
	require.Equal(t, "rejected", matchdomain.Rejected.String())
	require.Equal(t, matchdomain.Waiting, matchdomain.ParseUserStatus("waiting"))
	require.Equal(t, matchdomain.Matched, matchdomain.ParseUserStatus("matched"))
	require.Equal(t, matchdomain.Rejected, matchdomain.ParseUserStatus("rejected"))
	require.Equal(t, matchdomain.Waiting, matchdomain.ParseUserStatus("unknown"))
}

func TestUserStatusMarshalUnmarshalText(t *testing.T) {
	var s matchdomain.UserStatus = matchdomain.Matched
	b, err := s.MarshalText()
	require.NoError(t, err)
	require.Equal(t, []byte("matched"), b)

	var s2 matchdomain.UserStatus
	err = s2.UnmarshalText([]byte("rejected"))
	require.NoError(t, err)
	require.Equal(t, matchdomain.Rejected, s2)
}

func TestUserMarshalUnmarshalJSON(t *testing.T) {
	prefs := matchmaking.Preferences{MinAge: 20, MaxAge: 40, Gender: gender.Unspecified}
	u := matchdomain.NewUser("user3", 30, gender.Unspecified, prefs)
	u.SetStatus(matchdomain.Matched)

	data, err := json.Marshal(u)
	require.NoError(t, err)

	var u2 matchdomain.User
	err = json.Unmarshal(data, &u2)
	require.NoError(t, err)
	require.Equal(t, u.ID(), u2.ID())
	require.Equal(t, u.Age(), u2.Age())
	require.Equal(t, u.Gender(), u2.Gender())
	require.Equal(t, u.Preferences(), u2.Preferences())
	require.Equal(t, u.Status(), u2.Status())
}
