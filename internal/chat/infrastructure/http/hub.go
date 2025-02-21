package chathttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	chatcommands "github.com/xfrr/randomtalk/internal/chat/application/commands"
	chatqueries "github.com/xfrr/randomtalk/internal/chat/application/queries"
	chatconfig "github.com/xfrr/randomtalk/internal/chat/config"
)

// WebSocket Upgrader with proper settings
var upgrader = &websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins; customize for security
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents a WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan []byte // Buffered channel for outbound messages
	hub  *Hub
}

// Hub maintains active connections
type Hub struct {
	cfg        chatconfig.HubWebsocketServer
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	cmdBus     chatcommands.CommandBus
	queryBus   chatqueries.QueryBus
	logger     zerolog.Logger
}

type HubOption func(*Hub)

func WithLogger(logger zerolog.Logger) HubOption {
	return func(h *Hub) {
		h.logger = logger
	}
}

// NewHub creates a new WebSocket hub
func NewHub(cmdBus chatcommands.CommandBus, queryBus chatqueries.QueryBus, opts ...HubOption) *Hub {
	hub := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		cmdBus:     cmdBus,
		queryBus:   queryBus,
		logger:     zerolog.Nop(), // Avoid nil logger
		cfg: chatconfig.HubWebsocketServer{
			Address:             ":51500",
			Path:                "/",
			ReadBufferSize:      1024,
			ReadTimeoutSeconds:  10,
			WriteTimeoutSeconds: 10,
			IdleTimeoutSeconds:  10,
			WriteBufferSize:     1024,
			PongWaitSeconds:     60,
			PingPeriodSeconds:   54,
			MaxMessageSizeBytes: 1024,
			MaxConnections:      1000,
		},
	}

	for _, opt := range opts {
		opt(hub)
	}

	upgrader.ReadBufferSize = hub.cfg.ReadBufferSize
	upgrader.WriteBufferSize = hub.cfg.WriteBufferSize
	return hub
}

// Run starts the hub to manage clients
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)

		case client := <-h.unregister:
			h.removeClient(client)

		case msg := <-h.broadcast:
			h.broadcastMessage(msg)
		}
	}
}

func (h *Hub) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
	h.logger.Debug().Msg("client connected")
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client]; exists {
		delete(h.clients, client)
		close(client.send)
		h.logger.Debug().Msg("client disconnected")
	}
}

func (h *Hub) broadcastMessage(msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- msg:
		default:
			h.removeClient(client)
		}
	}
}

// Handle manages incoming WebSocket connections.
func (h *Hub) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to upgrade connection")
		return
	}

	if len(h.clients) >= h.cfg.MaxConnections {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseMessageTooBig, "server is full"))
		_ = conn.Close()
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
	}

	h.register <- client

	// Start client read and write pumps
	go client.readPump()
	go client.writePump()
}

// readPump reads messages from a client
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(int64(c.hub.cfg.MaxMessageSizeBytes))
	_ = c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.hub.cfg.ReadTimeoutSeconds) * time.Second))

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.hub.cfg.PongWaitSeconds) * time.Second))
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				c.hub.logger.Error().Err(err).Msg("unexpected close error")
			}
			break
		}

		var rawCommand struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}

		if err = json.Unmarshal(msg, &rawCommand); err != nil {
			respondError(c, "system", errors.New("invalid command"))
			continue
		}

		res, err := c.hub.cmdBus.Dispatch(context.Background(), rawCommand.Type, rawCommand.Payload)
		if err != nil {
			respondError(c, "system", err)
			continue
		}

		if res == nil {
			continue
		}

		responseEncoder, err := ServerResponseEncoders.GetEncoder(rawCommand.Type)
		if err != nil {
			c.hub.logger.Error().
				Str("command_type", rawCommand.Type).
				Err(err).
				Msg("failed to get encoder")
			continue
		}

		encodedResponse, err := responseEncoder(res)
		if err != nil {
			c.hub.logger.Error().
				Str("command_type", rawCommand.Type).
				Err(err).
				Msg("failed to encode response")
			continue
		}

		respond(c, encodedResponse)
	}
}

// writePump writes messages to a client
func (c *Client) writePump() {
	heartbeatTicker := time.NewTicker(time.Duration(c.hub.cfg.PingPeriodSeconds) * time.Second)
	defer func() {
		heartbeatTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.hub.cfg.WriteTimeoutSeconds) * time.Second)); err != nil {
				c.hub.logger.Error().Err(err).Msg("failed to set write deadline")
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.hub.logger.Error().Err(err).Msg("failed to get writer")
				return
			}

			_, err = w.Write(msg)
			if err != nil {
				c.hub.logger.Error().
					Stringer("client", c.conn.RemoteAddr()).
					Err(err).
					Msg("failed to write message")
				return
			}

			// Send queued messages in a single frame
			for range len(c.send) {
				_, _ = w.Write([]byte("\n"))
				_, _ = w.Write(<-c.send)
			}

			if err = w.Close(); err != nil {
				c.hub.logger.Error().Err(err).Msg("failed to close writer")
				return
			}

		case <-heartbeatTicker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(time.Duration(c.hub.cfg.WriteTimeoutSeconds) * time.Second)); err != nil {
				c.hub.logger.Error().Err(err).Msg("failed to set write deadline")
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
