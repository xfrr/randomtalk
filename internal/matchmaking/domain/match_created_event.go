package matchdomain

import (
	"github.com/xfrr/go-cqrsify/domain"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

// MatchCreatedEvent is an event that is published when a match  is created.
type MatchCreatedEvent struct {
	domain.BaseEvent

	MatchID                       string                  `json:"match_id"`
	MatchUserRequesterID          string                  `json:"match_user_requester_id"`
	MatchUserRequesterAge         int32                   `json:"match_user_requester_age"`
	MatchUserRequesterGender      gender.Gender           `json:"match_user_requester_gender"`
	MatchUserRequesterPreferences matchmaking.Preferences `json:"match_user_requester_preferences"`

	MatchUserMatchedID          string                  `json:"match_user_matched_id"`
	MatchUserMatchedAge         int32                   `json:"match_user_matched_age"`
	MatchUserMatchedGender      gender.Gender           `json:"match_user_matched_gender"`
	MatchUserMatchedPreferences matchmaking.Preferences `json:"match_user_matched_preferences"`
}

func (e MatchCreatedEvent) EventName() string {
	return "match_created"
}

func NewMatchCreatedEvent(
	match *Match,
	requesterUser User,
	matchedUser User,
) *MatchCreatedEvent {
	return &MatchCreatedEvent{
		BaseEvent:                     domain.NewEvent("match_created", domain.CreateEventAggregateRef(match)),
		MatchID:                       match.ID(),
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
