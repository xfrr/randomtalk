package chathttp

import (
	"errors"

	chatpb "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/chat/v1"
)

var (
	ErrNoEncoders     = errors.New("no encoders registered")
	ErrNoEncoderFound = errors.New("no encoder found for command")
)

var ServerResponseEncoders = &serverResponseEncoderRegistry{}

type ServerResponseEncoder func(any) (*chatpb.ServerMessage, error)

type serverResponseEncoderRegistry struct {
	encoders map[string]ServerResponseEncoder
}

func (r *serverResponseEncoderRegistry) Register(commandName string, encoder ServerResponseEncoder) {
	if r.encoders == nil {
		r.encoders = make(map[string]ServerResponseEncoder)
	}

	r.encoders[commandName] = encoder
}

func (r *serverResponseEncoderRegistry) GetEncoder(commandName string) (ServerResponseEncoder, error) {
	if r.encoders == nil {
		return nil, ErrNoEncoders
	}

	encoder, ok := r.encoders[commandName]
	if !ok {
		return nil, ErrNoEncoderFound
	}

	return encoder, nil
}
