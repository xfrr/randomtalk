package matchdomain

import (
	"fmt"
	"time"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

// Match is a domain entity representing a successful pairing of two users.
type Match struct {
	*aggregate.Base[string]

	requester *User
	match     *User
	createdAt time.Time
}

func (m *Match) ID() string {
	return m.AggregateID()
}

func (m *Match) Age() int {
	return m.requester.age
}

func (m *Match) Gender() gender.Gender {
	return m.requester.gender
}

func (m *Match) Preferences() matchmaking.MatchPreferences {
	if m.requester.preferences == nil {
		return *matchmaking.DefaultPreferences()
	}

	return *m.requester.preferences
}

func (m *Match) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Match) Requester() *User {
	return m.requester
}

func (m *Match) Candidate() *User {
	return m.match
}

func (m *Match) registerEventHandlers() {
	var matchCreatedEvent MatchCreatedEvent
	m.HandleEvent(matchCreatedEvent.EventName(), m.handleMatchCreatedEvent)
}

func (m *Match) handleMatchCreatedEvent(evt aggregate.Event) {
	payload, ok := evt.Payload().(*MatchCreatedEvent)
	if !ok {
		panic(fmt.Sprintf("unexpected event payload type: %T", evt.Payload()))
	}

	m.requester = &User{
		id:          payload.MatchUserRequesterID,
		age:         payload.MatchUserRequesterAge,
		gender:      payload.MatchUserRequesterGender,
		preferences: &payload.MatchUserRequesterPreferences,
	}

	m.match = &User{
		id:          payload.MatchUserMatchedID,
		age:         payload.MatchUserMatchedAge,
		gender:      payload.MatchUserMatchedGender,
		preferences: &payload.MatchUserMatchedPreferences,
	}

	m.createdAt = time.Now()
}

func (m *Match) validate() error {
	if m.ID() == "" {
		return ErrMatchIDNotProvided
	}

	if m.Requester() == nil {
		return ErrMatchRequesterNotProvided
	}

	if m.Candidate() == nil {
		return ErrMatchCandidateNotProvided
	}

	if m.Requester().ID() == m.Candidate().ID() {
		return ErrUserCannotMatchWithItself
	}
	return nil
}
