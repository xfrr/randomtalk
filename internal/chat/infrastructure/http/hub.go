package chathttp

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/xfrr/go-cqrsify/cqrs"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/xfrr/randomtalk/internal/chat/infrastructure/auth"
	"github.com/xfrr/randomtalk/internal/shared/messaging"
	"github.com/xfrr/randomtalk/internal/shared/semantic"

	chatcommands "github.com/xfrr/randomtalk/internal/chat/application/commands"
	chatqueries "github.com/xfrr/randomtalk/internal/chat/application/queries"
	chatconfig "github.com/xfrr/randomtalk/internal/chat/config"
	httpencoding "github.com/xfrr/randomtalk/internal/chat/infrastructure/http/encoding"
	chatpbv1 "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/chat/v1"
)

// NotificationConsumer is a component that consumes notifications.
type NotificationConsumer interface {
	Consume(ctx context.Context, notificationHandler func(ctx context.Context, notification *messaging.Event)) error
}

// WebSocket Upgrader with proper settings
var upgrader = &websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins; customize for security
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents a WebSocket connection
type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte // Buffered channel for outbound messages
	hub  *Hub
}

// Hub maintains active connections
type Hub struct {
	cfg                   *chatconfig.HubWebsocketServer
	clients               map[*Client]bool
	broadcast             chan []byte
	register              chan *Client
	unregister            chan *Client
	mu                    sync.RWMutex
	cmdBus                chatcommands.CommandBus
	queryBus              chatqueries.QueryBus
	notificationsConsumer NotificationConsumer
	logger                zerolog.Logger
}

type HubOption func(*Hub)

// WithLogger sets the logger for the hub instance.
func WithLogger(logger zerolog.Logger) HubOption {
	return func(h *Hub) {
		h.logger = logger
	}
}

// WithConfig sets the configuration for the hub instance.
func WithConfig(cfg *chatconfig.HubWebsocketServer) HubOption {
	return func(h *Hub) {
		h.cfg = cfg
	}
}

// NewHub creates a new WebSocket hub instance.
func NewHub(cmdBus chatcommands.CommandBus, queryBus chatqueries.QueryBus, notificationsConsumer NotificationConsumer, opts ...HubOption) *Hub {
	hub := &Hub{
		broadcast:             make(chan []byte),
		register:              make(chan *Client),
		unregister:            make(chan *Client),
		clients:               make(map[*Client]bool),
		cmdBus:                cmdBus,
		queryBus:              queryBus,
		logger:                zerolog.Nop(), // Avoid nil logger
		notificationsConsumer: notificationsConsumer,
		cfg: &chatconfig.HubWebsocketServer{
			Address: ":51000",
			Path:    "/ws",
			// 4096 is a common buffer size that reduces overhead for mid-to-large messages.
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			// A read timeout ensures you don't block indefinitely on slow or unresponsive clients.
			// 60s is a typical upper bound for most chat/real-time apps.
			ReadTimeoutSeconds: 60,
			// A write timeout of around 15-30 seconds is often enough unless you regularly send
			// large payloads or expect clients on very slow connections.
			WriteTimeoutSeconds: 15,
			// Idle timeout is how long the connection can stay open with no activity. Often,
			// production apps set this higher than 10s (e.g., 60-300s) so that legitimate
			// clients aren't disconnected too quickly.
			IdleTimeoutSeconds: 120,
			// The server should expect pongs within, say, 60s. The ping period should be
			// slightly less than the pong wait to avoid false timeouts.
			PongWaitSeconds:   60,
			PingPeriodSeconds: 54,
			// Limiting max message size protects against malicious or accidental large payloads.
			// 1 KiB might be too limiting for typical production scenarios; 64 KiB is more common.
			MaxMessageSizeBytes: 64 * 1024, // 65536 bytes
			// Enforcing a max connection count helps protect server resources. For high-scale
			// deployments, 1000 might be too low. Depending on your hardware and load tests,
			// you might allow 10k+ or handle scaling horizontally.
			MaxConnections: 10000,
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
func (h *Hub) Run(ctx context.Context) {
	go h.startNotificationsConsumer(ctx)

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
	h.logger.Debug().
		Str("client.address", client.conn.RemoteAddr().String()).
		Str("client.id", client.id).
		Msg("client registered")
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client]; exists {
		delete(h.clients, client)
		close(client.send)
		h.logger.Debug().
			Str("client.address", client.conn.RemoteAddr().String()).
			Str("client.id", client.id).
			Msg("client unregistered")
	}
}

// TODO: create an index for userID to client
// getClientByUserID retrieves a client by user ID
func (h *Hub) getClientByUserID(userID string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.id == userID {
			return client
		}
	}

	return nil
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

func (h *Hub) startNotificationsConsumer(ctx context.Context) {
	err := h.notificationsConsumer.Consume(ctx, func(ctx context.Context, notification *messaging.Event) {
		h.logger.Debug().Msg("received notification, sending to clients")

		var dataMap map[string]any
		if err := notification.DataAs(&dataMap); err != nil {
			h.logger.Error().Err(err).Msg("failed to decode notification data")
			notification.Reject()
			return
		}

		requesterUserID, ok := dataMap["match_user_requester_id"].(string)
		if !ok {
			h.logger.Error().Msg("failed to get requester user ID from notification")
			notification.Reject()
			return
		}

		matchedUserID, ok := dataMap["match_user_matched_id"].(string)
		if !ok {
			h.logger.Error().Msg("failed to get matched user ID from notification")
			notification.Reject()
			return
		}

		// send notification to requester
		firstUser := h.getClientByUserID(requesterUserID)
		if firstUser == nil {
			h.logger.Error().
				Str("requester_user_id", requesterUserID).
				Str("matched_user_id", matchedUserID).
				Msg("requester user not found")
			notification.Reject()
			return
		}

		// send notification to matched user
		secondUser := h.getClientByUserID(matchedUserID)
		if secondUser == nil {
			h.logger.Error().
				Str("requester_user_id", requesterUserID).
				Str("matched_user_id", matchedUserID).
				Msg("matched user not found")
			notification.Reject()
			return
		}

		// Create a new notification payload
		notificationDataMap := map[string]any{
			"user_requester_id": requesterUserID,
			"user_matched_id":   matchedUserID,
		}

		nprotoPayload, err := structpb.NewStruct(notificationDataMap)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to create struct from notification data")
			notification.Reject()
			return
		}

		nproto := &chatpbv1.ServerMessage{
			Kind: chatpbv1.Kind_KIND_SYSTEM,
			Data: &chatpbv1.ServerMessage_Notification{
				Notification: &chatpbv1.NotificationMessage{
					Type:    chatpbv1.NotificationMessage_TYPE_NEW_MATCH,
					Payload: nprotoPayload,
				},
			},
		}

		data, err := protojson.Marshal(nproto)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to marshal notification")
			notification.Reject()
			return
		}

		// send notification to both users
		firstUser.send <- data
		secondUser.send <- data

		h.logger.Debug().Msg("notification sent to users")
		// Acknowledge the notification
		notification.Ack()
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to start notifications consumer")
	}
}

// Handle manages incoming WebSocket connections.
func (h *Hub) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to upgrade connection")
		return
	}

	// get content-type from request
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	if len(h.clients) >= h.cfg.MaxConnections {
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseMessageTooBig, "server is full"))
		_ = conn.Close()
		return
	}

	client := &Client{
		id:   uuid.New().String(),
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
	}

	h.register <- client

	// Start client read and write pumps
	go client.readPump(contentType)
	go client.writePump()
}

// readPump reads messages from a client
func (c *Client) readPump(contentType string) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(int64(c.hub.cfg.MaxMessageSizeBytes))
	_ = c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.hub.cfg.ReadTimeoutSeconds) * time.Second))

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.hub.cfg.PongWaitSeconds) * time.Second))
	})

	logger := c.hub.logger.With().
		Str("client", c.conn.RemoteAddr().String()).
		Logger()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil && websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			logger.Error().Err(err).Msg("failed to read message")
			break
		} else if err != nil {
			break
		}

		ctx := context.Background()

		logger = logger.With().
			Ctx(ctx).
			Str(semantic.EnduserIDKey, c.id).
			Str(semantic.ClientAddressKey, c.conn.RemoteAddr().String()).
			Int(semantic.HTTPRequestBodySizeKey, len(msg)).
			Str(semantic.MessageTypeKey, "command").
			Str(semantic.HTTPRequestContentTypeKey, contentType).
			Logger()

		cmd, err := httpencoding.DecodeCommand(contentType, msg)
		if err != nil {
			logger.Error().Err(err).Msg("failed to decode message")
			respondError(c, "system", err)
			continue
		}

		ctx = auth.ContextWithUserID(ctx, c.id)
		res, err := cqrs.Dispatch[any](ctx, c.hub.cmdBus, cmd)
		if err != nil {
			logger.Error().Err(err).Msg("failed to dispatch command")
			respondError(c, "system", err)
			continue
		}

		if res == nil {
			continue
		}
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
