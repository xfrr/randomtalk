package httpencoding

import (
	"encoding/json"
	"io"
)

// ApplicationJSON is the application/json MIME type.
const ApplicationJSON = "application/json"

// jsonEncoder is a generic JSON Encoder for type T.
type jsonEncoder[T any] struct{}

// JSONEncoder instantiates a new JSONEncoder.
func JSONEncoder[T any]() *jsonEncoder[T] {
	return &jsonEncoder[T]{}
}

// Encode implements the Encoder[T] interface using JSON.
func (je *jsonEncoder[T]) Encode(w io.Writer, cmd T) error {
	return json.NewEncoder(w).Encode(cmd)
}

// JSONDecoder is a generic JSON Decoder for type T.
type JSONDecoder[T any] struct{}

// NewJSONDecoder instantiates a new JSONDecoder.
func NewJSONDecoder[T any]() *JSONDecoder[T] {
	return &JSONDecoder[T]{}
}

// Decode implements the Decoder[T] interface using JSON.
func (jd *JSONDecoder[T]) Decode(r io.Reader) (T, error) {
	var result T
	err := json.NewDecoder(r).Decode(&result)
	return result, err
}
