package chatdomain

import (
	"github.com/xfrr/go-cqrsify/aggregate"
	chatdomaineventsv1 "github.com/xfrr/randomtalk/internal/chat/domain/events/v1"
	"github.com/xfrr/randomtalk/internal/shared/gender"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"
)

// chatSessionCreatedDomainEventHandler is a domain event handler for ChatSessionCreated events.
func (cs *ChatSession) chatSessionCreatedDomainEventHandlerV1(evt aggregate.Event) {
	payload, _ := evt.Payload().(chatdomaineventsv1.ChatSessionCreated)

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
}
