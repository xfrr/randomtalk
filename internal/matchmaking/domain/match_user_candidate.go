package matchdomain

import "sync"

// MatchUserCandidate represents a user candidate for matchmaking.
type MatchUserCandidate struct {
	User

	acceptedCh   chan struct{}
	acceptedOnce sync.Once

	rejectedCh   chan struct{}
	rejectedOnce sync.Once

	abortedCh   chan error
	abortedOnce sync.Once
}

// NewMatchUserCandidate creates a new MatchUserCandidate.
func NewMatchUserCandidate(user User) *MatchUserCandidate {
	return &MatchUserCandidate{
		User:       user,
		acceptedCh: make(chan struct{}, 1),
		rejectedCh: make(chan struct{}, 1),
		abortedCh:  make(chan error, 1),
	}
}

// Accept accepts the match candidate only once.
func (m *MatchUserCandidate) Accept() {
	m.acceptedOnce.Do(func() {
		close(m.acceptedCh)
	})
}

// Reject rejects the match candidate only once.
func (m *MatchUserCandidate) Reject() {
	m.rejectedOnce.Do(func() {
		close(m.rejectedCh)
	})
}

// Abort aborts the match candidate only once.
func (m *MatchUserCandidate) Abort(err error) {
	m.abortedOnce.Do(func() {
		m.abortedCh <- err
	})
}

// Accepted returns the accepted channel.
func (m *MatchUserCandidate) Accepted() <-chan struct{} {
	return m.acceptedCh
}

// Rejected returns the rejected channel.
func (m *MatchUserCandidate) Rejected() <-chan struct{} {
	return m.rejectedCh
}

// Aborted returns the aborted channel.
func (m *MatchUserCandidate) Aborted() <-chan error {
	return m.abortedCh
}
