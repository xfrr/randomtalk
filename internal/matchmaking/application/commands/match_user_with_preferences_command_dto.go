package matchcommands

import (
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

type MatchUserWithPreferencesCommand struct {
	UserID               string                       `json:"user_id"`
	UserAge              int                          `json:"user_age"`
	UserGender           gender.Gender                `json:"user_gender"`
	UserMatchPreferences matchmaking.MatchPreferences `json:"user_match_preferences"`
}

func (c MatchUserWithPreferencesCommand) CommandName() string {
	return "match_user_with_preferences"
}

type MatchUserWithPreferencesResponse struct {
	MatchID string `json:"match_id"`
}
