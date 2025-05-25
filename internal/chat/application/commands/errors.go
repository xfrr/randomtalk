package chatcommands

import "errors"

var (
	ErrMissingUserIDFromContext = errors.New("missing user ID from context")
)
