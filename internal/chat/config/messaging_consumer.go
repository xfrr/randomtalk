package chatconfig

// MessagingEngine is the type of the consumer engine.
type MessagingEngine string

func (t MessagingEngine) String() string {
	return string(t)
}

func (t MessagingEngine) IsValid() bool {
	switch t {
	case MessagingEngineNATS:
		return true
	default:
		return false
	}
}

const (
	// MessagingEngineNATS is the NATS engine.
	MessagingEngineNATS MessagingEngine = "nats"
)

type MatchNotificationsConsumerConfig struct {
	Engine     MessagingEngine `env:"ENGINE" default:"nats"`
	Name       string          `env:"NAME" default:"randomtalk_chat_match_events_consumer"`
	StreamName string          `env:"STREAM_NAME" default:"randomtalk_matchmaking_match_events"`
}
