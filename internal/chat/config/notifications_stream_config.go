package chatconfig

// NotificationStreamConfig holds the configuration for the JetStream stream used for notifications.
type NotificationStreamConfig struct {
	Name string `envconfig:"NOTIFICATION_STREAM_NAME" default:"randomtalk_chat_notifications"`
}
