package matchcommands

import (
	"github.com/xfrr/randomtalk/pkg/gender"

	matchdomain "github.com/xfrr/randomtalk/internal/matchmaking/domain"
)

type MatchUserWithPreferencesCommand struct {
	UserID               string                       `json:"user_id"`
	UserAge              int                          `json:"user_age"`
	UserGender           gender.Gender                `json:"user_gender"`
	UserMatchPreferences matchdomain.MatchPreferences `json:"user_match_preferences"`
}

func (c MatchUserWithPreferencesCommand) CommandName() string {
	return "match_user_with_preferences"
}

type MatchUserWithPreferencesResponse struct {
	MatchID string `json:"match_id"`
}
