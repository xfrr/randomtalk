package chatdomain

import (
	"fmt"

	"github.com/xfrr/go-cqrsify/domain"
	chatdomaineventsv1 "github.com/xfrr/randomtalk/internal/chat/domain/events/v1"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

// chatSessionCreatedDomainEventHandler is a domain event handler for ChatSessionCreated events.
func (cs *ChatSession) chatSessionCreatedDomainEventHandlerV1(evt domain.Event) error {
	payload, ok := evt.(chatdomaineventsv1.ChatSessionCreated)
	if !ok {
		return fmt.Errorf("unexpected event type: %T, expected: %T", evt, chatdomaineventsv1.ChatSessionCreated{})
	}

	if cs.state == nil {
		cs.state = &chatSessionState{
			User: new(User),
		}
	}

	cs.state.User = &User{
		id:       ID(payload.UserID),
		nickname: payload.UserNickname,
		age:      payload.UserAge,
		gender:   gender.Parse(payload.UserGender),
		matchPreferences: matchmaking.
			DefaultPreferences().
			WithGender(gender.Parse(payload.UserPreference.Gender)).
			WithMinAge(payload.UserPreference.MinAge).
			WithMaxAge(payload.UserPreference.MaxAge).
			WithInterests(payload.UserPreference.Interests),
	}
	return nil
}
