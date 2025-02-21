package domain_error

import "errors"

// Error represents a domain error.
type Error struct {
	err error
	msg string

	aggregateID      string
	aggregateName    string
	aggregateVersion int
}

// New creates a new domain error.
func New(msg string) *Error {
	return &Error{
		msg: msg,
	}
}

// Error returns the error message.
func (e *Error) Error() string {
	var wrapped string
	if e.err != nil {
		wrapped = e.err.Error()
	}
	if wrapped != "" {
		return e.msg + ": " + wrapped
	}
	return e.msg
}

// AggregateID returns the aggregate ID.
func (e *Error) AggregateID() string {
	return e.aggregateID
}

// AggregateName returns the aggregate name.
func (e *Error) AggregateName() string {
	return e.aggregateName
}

// AggregateVersion returns the aggregate version.
func (e *Error) AggregateVersion() int {
	return e.aggregateVersion
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.err
}

// WithAggregateID sets the aggregate ID.
func (e *Error) WithAggregateID(id string) *Error {
	e.aggregateID = id
	return e
}

// WithAggregateName sets the aggregate name.
func (e *Error) WithAggregateName(name string) *Error {
	e.aggregateName = name
	return e
}

// WithAggregateVersion sets the aggregate version.
func (e *Error) WithAggregateVersion(version int) *Error {
	e.aggregateVersion = version
	return e
}

// Wrap sets the error.
func (e *Error) Wrap(err error) *Error {
	e.err = err
	return e
}

// IsDomainError checks if the error is a domain error.
func IsDomainError(err error) bool {
	var domainErr *Error
	return errors.As(err, &domainErr)
}
