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
	_, errMsg := formatErrorResponse(err)
	rawMsg := &chatpb.ServerMessage{
		Id:        msgID,
		Type:      chatpb.ServerMessage_TYPE_ERROR,
		Timestamp: timestamppb.Now(),
		Payload: &chatpb.ServerMessage_Payload{
			Content: &chatpb.ServerMessage_Payload_Error{
				Error: status.New(codes.Internal, errMsg).Proto(),
			},
		},
	}
	respond(client, rawMsg)
}

func respond(client *Client, msg *chatpb.ServerMessage) {
	rawMsg, _ := protojson.Marshal(msg)
	client.send <- rawMsg
}

func formatErrorResponse(err error) (int32, string) {
	var code int32
	var message string

	switch {
	case errors.Is(err, cqrs.ErrHandlerNotFound):
		code = http.StatusNotFound
		message = "Command not found"
	default:
		code = http.StatusInternalServerError
		message = "Internal server error"
	}
	return code, message
}
