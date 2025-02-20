package matchmakingconfig

import "github.com/kelseyhightower/envconfig"

// ConsumerEngine is the type of the consumer engine.
type ConsumerEngine string

func (t ConsumerEngine) String() string {
	return string(t)
}

func (t ConsumerEngine) IsValid() bool {
	switch t {
	case ConsumerEngineNATS:
		return true
	default:
		return false
	}
}

const (
	// ConsumerEngineNATS is the NATS engine.
	ConsumerEngineNATS ConsumerEngine = "nats"
)

// ChatNotificationsConsumerConfig holds the configuration for the chat notifications consumer.
type ChatNotificationsConsumerConfig struct {
	Engine     ConsumerEngine `envconfig:"CHAT_NOTIFICATIONS_CONSUMER_ENGINE" default:"nats"`
	Name       string         `envconfig:"CHAT_NOTIFICATIONS_CONSUMER_NAME" default:"randomtalk_matchmaking_chat_notifications_consumer"`
	StreamName string         `envconfig:"CHAT_NOTIFICATIONS_STREAM_NAME" default:"randomtalk_chat_notifications"`
}

func mustLoadChatNotificationsConsumerConfig() ChatNotificationsConsumerConfig {
	var cfg ChatNotificationsConsumerConfig
	envconfig.MustProcess(envPrefix, &cfg)
	return cfg
}
