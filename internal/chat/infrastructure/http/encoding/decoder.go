package httpencoding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/xfrr/go-cqrsify/messaging"
	chatcommands "github.com/xfrr/randomtalk/internal/chat/application/commands"
)

var (
	decoders = map[string]map[string]Decoder[any]{
		chatcommands.CreateChatSessionCommandType: {
			ApplicationJSON: AnyDecoder(NewJSONDecoder[chatcommands.CreateChatSessionCommand]()),
		},
	}
)

// Decoder decodes a given message to a given T type.
type Decoder[T any] interface {
	Decode(r io.Reader) (T, error)
}

type anyDecoder[T any] struct {
	decoder Decoder[T]
}

func (d anyDecoder[T]) Decode(r io.Reader) (any, error) {
	return d.decoder.Decode(r)
}

// DecodeCommand decodes a command from a message using the given content type.
func DecodeCommand(contentType string, msg []byte) (messaging.Command, error) {
	var rawCommand struct {
		Kind string `json:"kind"`
		Data struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		} `json:"data"`
	}

	if err := json.Unmarshal(msg, &rawCommand); err != nil {
		return nil, err
	}

	switch rawCommand.Kind {
	case "command":
		decoder, ok := decoders[rawCommand.Data.Type][contentType]
		if !ok {
			return nil, fmt.Errorf("unsupported command type %s with content type %s", rawCommand.Data.Type, contentType)
		}

		byteReader := bytes.NewReader(rawCommand.Data.Payload)
		cmd, err := decoder.Decode(byteReader)
		if err != nil {
			return nil, err
		}

		castedCmd, ok := cmd.(messaging.Command)
		if !ok {
			return nil, errors.New("decoded command is not a messaging.Command")
		}
		return castedCmd, nil
		// TODO: handle other message kinds
	default:
		return nil, fmt.Errorf("unsupported message kind %s", rawCommand.Kind)
	}
}
