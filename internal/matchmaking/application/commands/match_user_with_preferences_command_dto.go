package matchcommands

import (
	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

type MatchUserWithPreferencesCommand struct {
	messaging.BaseCommand

	UserID          string                  `json:"user_id"`
	UserAge         int32                   `json:"user_age"`
	UserGender      gender.Gender           `json:"user_gender"`
	UserPreferences matchmaking.Preferences `json:"user_match_preferences"`
}

type MatchUserWithPreferencesResponse struct {
	MatchID string `json:"match_id"`
}
