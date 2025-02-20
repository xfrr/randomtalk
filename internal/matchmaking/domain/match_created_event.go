package matchdomain

import (
	"github.com/xfrr/randomtalk/pkg/gender"
)

// MatchCreatedEvent is an event that is published when a match  is created.
type MatchCreatedEvent struct {
	MatchID                       string           `json:"match_id"`
	MatchUserRequesterID          string           `json:"match_user_requester_id"`
	MatchUserRequesterAge         int              `json:"match_user_requester_age"`
	MatchUserRequesterGender      gender.Gender    `json:"match_user_requester_gender"`
	MatchUserRequesterPreferences MatchPreferences `json:"match_user_requester_preferences"`

	MatchUserMatchedID          string           `json:"match_user_matched_id"`
	MatchUserMatchedAge         int              `json:"match_user_matched_age"`
	MatchUserMatchedGender      gender.Gender    `json:"match_user_matched_gender"`
	MatchUserMatchedPreferences MatchPreferences `json:"match_user_matched_preferences"`
}

func (e MatchCreatedEvent) EventName() string {
	return "created"
}

func NewMatchCreatedEvent(
	matchID string,
	requesterUser User,
	matchedUser User,
) *MatchCreatedEvent {
	return &MatchCreatedEvent{
		MatchID:                       matchID,
		MatchUserRequesterID:          requesterUser.ID(),
		MatchUserRequesterAge:         requesterUser.Age(),
		MatchUserRequesterGender:      requesterUser.Gender(),
		MatchUserRequesterPreferences: requesterUser.Preferences(),
		MatchUserMatchedID:            matchedUser.ID(),
		MatchUserMatchedAge:           matchedUser.Age(),
		MatchUserMatchedGender:        matchedUser.Gender(),
		MatchUserMatchedPreferences:   matchedUser.Preferences(),
	}
}
