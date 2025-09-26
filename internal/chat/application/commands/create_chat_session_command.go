package chatcommands

import "github.com/xfrr/go-cqrsify/messaging"

type CreateChatSessionCommand struct {
	messaging.BaseCommand
	CommandInfo

	UserNickname                 string   `json:"user_nickname"`
	UserAge                      int32    `json:"user_age"`
	UserGender                   string   `json:"user_gender"`
	UserMatchPreferenceMinAge    int32    `json:"user_match_preference_min_age"`
	UserMatchPreferenceMaxAge    int32    `json:"user_match_preference_max_age"`
	UserMatchPreferenceGender    string   `json:"user_match_preference_gender"`
	UserMatchPreferenceInterests []string `json:"user_match_preference_interests"`
}

type CreateChatSessionResponse struct {
	ChatSessionID string `json:"chat_session_id"`
}
