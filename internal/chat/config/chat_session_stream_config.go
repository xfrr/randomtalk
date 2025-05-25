package chatconfig

// ChatSessionStreamConfig holds the configuration for the JetStream stream used for notifications.
type ChatSessionStreamConfig struct {
	Name string `env:"CHAT_SESSION_STREAM_NAME" default:"randomtalk_chat_sessions"`
}
