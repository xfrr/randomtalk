package chatdomaineventsv1

import "github.com/xfrr/go-cqrsify/domain"

// ChatSessionCreated is an event that is published when a match  is created.
type ChatSessionCreated struct {
	domain.BaseEvent

	SessionID      string   `json:"session_id"`
	UserID         string   `json:"user_id"`
	UserNickname   string   `json:"user_nickname"`
	UserAge        int32    `json:"user_age"`
	UserGender     string   `json:"user_gender"`
	UserPreference UserPref `json:"user_preference"`
}

// UserPref is a struct that holds the user's preferences.
type UserPref struct {
	MinAge    int32    `json:"min_age"`
	MaxAge    int32    `json:"max_age"`
	Gender    string   `json:"gender"`
	Interests []string `json:"interests"`
}

func (e ChatSessionCreated) EventName() string {
	return "chat_session_created"
}
