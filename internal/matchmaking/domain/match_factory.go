package matchdomain

import (
	"github.com/google/uuid"
	"github.com/xfrr/go-cqrsify/aggregate"

	domainerror "github.com/xfrr/randomtalk/internal/shared/domain-error"
)

const EventSourceName = "randomtalk.matchmaking"
const MatchAggregateName = "match"

var (
	// ErrMatchIDNotProvided is returned when the match  ID is not provided.
	ErrMatchIDNotProvided = domainerror.New("match ID not provided")

	// ErrMatchRequesterNotProvided is returned when the match  requester is not provided.
	ErrMatchRequesterNotProvided = domainerror.New("match requester not provided")

	// ErrMatchCandidateNotProvided is returned when the match  candidate is not provided.
	ErrMatchCandidateNotProvided = domainerror.New("match candidate not provided")

	// ErrMatchRequesterPreferencesNotProvided is returned when the match  requester preferences are not provided.
	ErrMatchRequesterPreferencesNotProvided = domainerror.New("match requester preferences not provided")

	// ErrMatchCandidatePreferencesNotProvided is returned when the match  candidate preferences are not provided.
	ErrMatchCandidatePreferencesNotProvided = domainerror.New("match candidate preferences not provided")

	// ErrUserCannotMatchWithItself is returned when the match  requester tries to match with itself.
	ErrUserCannotMatchWithItself = domainerror.New("user cannot match with itself")
)

func NewMatch(
	msid MatchID,
	requesterUser User,
	matchedUser User,
) (*Match, error) {
	match := newMatch(msid)

	event := NewMatchCreatedEvent(
		msid.String(),
		requesterUser,
		matchedUser,
	)

	err := aggregate.RaiseEvent(
		match,
		uuid.New().String(),
		event.EventName(),
		event,
	)
	if err != nil {
		return nil, err
	}

	if validateErr := match.validate(); validateErr != nil {
		return nil, validateErr
	}

	return match, nil
}

func NewMatchFromEvents(id MatchID, events ...aggregate.Event) (*Match, error) {
	match := newMatch(id)

	err := aggregate.RestoreStateFromHistory(match, events)
	if err != nil {
		return nil, err
	}

	if validateErr := match.validate(); validateErr != nil {
		return nil, validateErr
	}

	return match, nil
}

func newMatch(id MatchID) *Match {
	match := &Match{
		Base: aggregate.New(string(id), MatchAggregateName),
	}

	match.registerEventHandlers()
	return match
}
