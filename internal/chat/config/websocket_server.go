package chatconfig

// HubWebsocketServer
type HubWebsocketServer struct {
	// Address is the address the websocket server will listen on.
	Address string `env:"ADDRESS" default:":51000"`
	// Path is the path the websocket server will listen on.
	Path string `env:"PATH" default:"/sessions"`
	// ReadBufferSize is the size of the read buffer for the websocket connection.
	ReadBufferSize int `env:"READ_BUFFER_SIZE" default:"1024"`
	// ReadTimeoutSeconds is the maximum time the server will wait for a read from the client.
	ReadTimeoutSeconds int `env:"READ_TIMEOUT_SECONDS" default:"60"`
	// WriteTimeoutSeconds is the maximum time the server will wait for a write to the client.
	WriteTimeoutSeconds int `env:"WRITE_TIMEOUT_SECONDS" default:"60"`
	// IdleTimeoutSeconds is the maximum time the server will wait for a new request.
	IdleTimeoutSeconds int `env:"IDLE_TIMEOUT_SECONDS" default:"60"`
	// WriteBufferSize is the size of the write buffer for the websocket connection.
	WriteBufferSize int `env:"WRITE_BUFFER_SIZE" default:"1024"`
	// PongWaitSeconds is the maximum time the server will wait for a ping from the client.
	PongWaitSeconds int `env:"PONG_WAIT" default:"60"`
	// PingPeriodSeconds is the time between pings from the server to the client.
	PingPeriodSeconds int `env:"PING_PERIOD" default:"54"`
	// MaxMessageSizeBytes is the maximum size of a message that can be read from the client.
	MaxMessageSizeBytes int `env:"MAX_MESSAGE_SIZE" default:"1024"`
	// MaxConnections is the maximum number of connections the server will accept.
	MaxConnections int `env:"MAX_CONNECTIONS" default:"1000"`
}
