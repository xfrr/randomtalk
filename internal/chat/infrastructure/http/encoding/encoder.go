package httpencoding

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	chatpbv1 "github.com/xfrr/randomtalk/proto/gen/go/randomtalk/chat/v1"
)

var encoders = map[string]map[string]Encoder[any]{
	reflect.TypeOf(&chatpbv1.ServerMessage{}).Name(): {
		ApplicationJSON: AnyEncoder(JSONEncoder[*chatpbv1.ServerMessage]()),
	},
}

// Encoder encodes a given message to a given T type.
type Encoder[T any] interface {
	Encode(w io.Writer, msg T) error
}

// anyEncoder is a generic encoder for any type.
type anyEncoder[T any] struct {
	encoder Encoder[T]
}

// Encode encodes a message to a given type.
func (e anyEncoder[T]) Encode(w io.Writer, msg any) error {
	castedMsg, ok := msg.(T)
	if !ok {
		return fmt.Errorf("cannot encode message of type %T as %T", msg, castedMsg)
	}

	return e.encoder.Encode(w, castedMsg)
}

// EncodeMessage encodes the given message to the specified content type.
// It returns the encoded bytes or an error if the encoding fails.
// If the content type is not supported, it returns ErrNotificationEncoderNotFound.
func EncodeMessage(contentType string, msg any) ([]byte, error) {
	encoder, ok := encoders[reflect.TypeOf(msg).Name()][contentType]
	if !ok {
		return nil, ErrNotificationEncoderNotFound
	}

	var buf bytes.Buffer
	if err := encoder.Encode(&buf, msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
