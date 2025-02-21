package chatconfig

const envPrefix = "RANDOMTALK_CHAT"

// HubWebsocketServer
type HubWebsocketServer struct {
	// Address is the address the websocket server will listen on.
	Address string `envconfig:"WEBSOCKET_SERVER_ADDRESS" default:":51500"`
	// Path is the path the websocket server will listen on.
	Path string `envconfig:"WEBSOCKET_SERVER_PATH" default:"/sessions"`
	// ReadBufferSize is the size of the read buffer for the websocket connection.
	ReadBufferSize int `envconfig:"WEBSOCKET_SERVER_READ_BUFFER_SIZE" default:"1024"`
	// ReadTimeoutSeconds is the maximum time the server will wait for a read from the client.
	ReadTimeoutSeconds int `envconfig:"WEBSOCKET_SERVER_READ_TIMEOUT_SECONDS" default:"10"`
	// WriteTimeoutSeconds is the maximum time the server will wait for a write to the client.
	WriteTimeoutSeconds int `envconfig:"WEBSOCKET_SERVER_WRITE_TIMEOUT_SECONDS" default:"10"`
	// IdleTimeoutSeconds is the maximum time the server will wait for a new request.
	IdleTimeoutSeconds int `envconfig:"WEBSOCKET_SERVER_IDLE_TIMEOUT_SECONDS" default:"10"`
	// WriteBufferSize is the size of the write buffer for the websocket connection.
	WriteBufferSize int `envconfig:"WEBSOCKET_SERVER_WRITE_BUFFER_SIZE" default:"1024"`
	// PongWaitSeconds is the maximum time the server will wait for a ping from the client.
	PongWaitSeconds int `envconfig:"PONG_WAIT" default:"60"`
	// PingPeriodSeconds is the time between pings from the server to the client.
	PingPeriodSeconds int `envconfig:"PING_PERIOD" default:"54"`
	// MaxMessageSizeBytes is the maximum size of a message that can be read from the client.
	MaxMessageSizeBytes int `envconfig:"MAX_MESSAGE_SIZE" default:"1024"`
	// MaxConnections is the maximum number of connections the server will accept.
	MaxConnections int `envconfig:"MAX_CONNECTIONS" default:"1000"`
}
