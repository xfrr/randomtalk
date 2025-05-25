package httpencoding

// AnyDecoder is a generic decoder for any type.
func AnyDecoder[T any](decoder Decoder[T]) Decoder[any] {
	return anyDecoder[T]{decoder}
}

// AnyEncoder is a generic encoder for any type.
func AnyEncoder[T any](encoder Encoder[T]) Encoder[any] {
	return anyEncoder[T]{encoder}
}
