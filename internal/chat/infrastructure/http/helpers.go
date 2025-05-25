package chathttp

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/xfrr/go-cqrsify/cqrs"
	chatconfig "github.com/xfrr/randomtalk/internal/chat/config"
	chatpb "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/chat/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ServeHTTP(cfg chatconfig.HubWebsocketServer, handler http.HandlerFunc) error {
	http.Handle(cfg.Path, handler)

	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler:        handler,
		ReadTimeout:    time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:    time.Duration(cfg.IdleTimeoutSeconds) * time.Second,
		MaxHeaderBytes: 1 << 10, // 1 KB
	}

	return server.Serve(lis)
}

func respondError(client *Client, msgID string, err error) {
	code, errMsg := formatErrorResponse(err)
	rawMsg := &chatpb.ServerMessage{
		Kind: chatpb.Kind_KIND_SYSTEM,
		Data: &chatpb.ServerMessage_Error{
			Error: &chatpb.ErrorMessage{
				Status:    status.New(codes.Code(code), errMsg).Proto(),
				Timestamp: timestamppb.Now(),
			},
		},
	}
	respond(client, rawMsg)
}

func respond(client *Client, msg *chatpb.ServerMessage) {
	rawMsg, _ := protojson.Marshal(msg)
	client.send <- rawMsg
}

func formatErrorResponse(err error) (uint32, string) {
	var code uint32
	var message string

	switch {
	case errors.Is(err, cqrs.ErrHandlerNotFound):
		code = http.StatusNotFound
		message = "Command not found: " + err.Error()
	default:
		code = http.StatusInternalServerError
		message = "Internal server error: " + err.Error()
	}
	return code, message
}
