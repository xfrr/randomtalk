package httpencoding

import "errors"

// ErrNotificationEncoderNotFound is the error returned when the encoder for a notification is not found.
var ErrNotificationEncoderNotFound = errors.New("notification encoder not found")
