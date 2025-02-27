package identity

import "github.com/google/uuid"

// NewUUID() generates a new UUID.
func NewUUID() ID {
	return ID(uuid.New().String())
}
