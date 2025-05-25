package chatdomain

import (
	"fmt"

	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/randomtalk/internal/shared/identity"
	"github.com/xfrr/randomtalk/internal/shared/matchmaking"

	chatdomaineventsv1 "github.com/xfrr/randomtalk/internal/chat/domain/events/v1"
	domain_error "github.com/xfrr/randomtalk/internal/shared/domain"
)

const EventSourceName = "randomtalk.chat"
const AggregateName = "session"

var (
	ErrInvalidChatSessionID = domain_error.New("invalid chat session unique identifier")
)

type ID = identity.ID
type MatchPreferences = matchmaking.Preferences

// ChatSession is the main Aggregate of the Chat bounded context.
// It represents a chat session between two users.
type ChatSession struct {
	*aggregate.Base[string]

	state *chatSessionState
}

type chatSessionState struct {
	User *User
}

// NewChatSession creates a new ChatSession instance.
func NewChatSession(id ID, user User) (*ChatSession, error) {
	cs := newChatSession(id)

	err := cs.raiseChatSessionCreatedEvent(user)
	if err != nil {
		return &ChatSession{}, err
	}

	if err = cs.validate(); err != nil {
		return &ChatSession{}, err
	}

	return cs, nil
}

// NewChatSessionFromEvents creates a new ChatSession instance from a list of events.
func NewChatSessionFromEvents(id ID, events []aggregate.Event) (*ChatSession, error) {
	cs := newChatSession(id)

	err := aggregate.RestoreStateFromHistory(cs, events)
	if err != nil {
		return cs, fmt.Errorf("failed to restore ChatSession state from history: %w", err)
	}

	if validateErr := cs.validate(); validateErr != nil {
		return cs, validateErr
	}

	return cs, nil
}

func newChatSession(id ID) *ChatSession {
	cs := &ChatSession{
		Base: aggregate.New(id.String(), AggregateName),
	}

	cs.registerEventHandlers()
	return cs
}

// ID returns the ChatSession ID.
func (cs ChatSession) ID() ID {
	return ID(cs.AggregateID())
}

// Users returns the ChatSession users.
func (cs ChatSession) User() *User {
	if cs.state == nil {
		return nil
	}
	return cs.state.User
}

func (cs *ChatSession) raiseChatSessionCreatedEvent(user User) error {
	chatSessionCreatedEvent := chatdomaineventsv1.ChatSessionCreated{
		SessionID:    cs.ID().String(),
		UserID:       user.ID().String(),
		UserNickname: user.Nickname(),
		UserAge:      user.Age(),
		UserGender:   user.Gender().String(),
		UserPreference: chatdomaineventsv1.UserPref{
			MinAge:    user.MatchPreferences().MinAge,
			MaxAge:    user.MatchPreferences().MaxAge,
			Gender:    user.MatchPreferences().Gender.String(),
			Interests: user.MatchPreferences().Interests,
		},
	}

	err := aggregate.RaiseEvent(
		cs,
		identity.NewUUID().String(),
		chatSessionCreatedEvent.EventName(),
		chatSessionCreatedEvent,
	)
	if err != nil {
		return fmt.Errorf("failed to raise ChatSessionCreated event: %w", err)
	}

	return nil
}

func (cs ChatSession) validate() error {
	if cs.ID() == "" {
		return ErrInvalidChatSessionID
	}
	return nil
}

func (cs *ChatSession) registerEventHandlers() {
	var chatSessionCreatedEvent chatdomaineventsv1.ChatSessionCreated

	cs.HandleEvent(chatSessionCreatedEvent.EventName(), cs.chatSessionCreatedDomainEventHandlerV1)
}
