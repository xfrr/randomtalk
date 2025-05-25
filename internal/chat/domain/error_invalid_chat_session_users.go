package chatdomain

import (
	"fmt"

	domain_error "github.com/xfrr/randomtalk/internal/shared/domain"
)

var ErrInvalidChatSessionUsers = invalidChatSessionUsersError{
	reason: domain_error.New("invalid number of users"),
}

type invalidChatSessionUsersError struct {
	reason   *domain_error.Error
	expected int
	actual   int
}

func (e *invalidChatSessionUsersError) CountMismatch(expected, actual int) *invalidChatSessionUsersError {
	e.expected = expected
	e.actual = actual
	return e
}

func (e *invalidChatSessionUsersError) Error() string {
	return fmt.Sprintf("%s: expected %d, got %d", e.reason.Error(), e.expected, e.actual)
}
