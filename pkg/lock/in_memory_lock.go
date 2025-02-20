package lock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// In-Memory Lock Adapter
// =============================================================================

// inMemoryRecord holds the expiry time and the token that identifies the lock owner.
type inMemoryRecord struct {
	expiry time.Time
	token  string
}

// InMemoryLockAdapter is an in-memory implementation of LockAdapter.
// It is suitable for local testing or single-process scenarios. Because it does
// not coordinate across multiple servers, it is not recommended for real
// distributed environments.
type InMemoryLockAdapter struct {
	mu    sync.Mutex
	locks map[string]inMemoryRecord
}

// NewInMemoryLockAdapter creates a new instance of InMemoryLockAdapter.
func NewInMemoryLockAdapter() *InMemoryLockAdapter {
	return &InMemoryLockAdapter{
		locks: make(map[string]inMemoryRecord),
	}
}

// AcquireLock acquires an in-memory lock if not already held or if the current holder's
// lock has expired. It returns a unique token identifying the caller as lock owner.
func (a *InMemoryLockAdapter) AcquireLock(_ context.Context, key string, expiration time.Duration) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()

	// Check if lock is currently held and has not expired.
	rec, exists := a.locks[key]
	if exists && now.Before(rec.expiry) {
		return "", errors.New("lock is already held by another token")
	}

	// Generate a new token (random).
	token, err := generateRandomToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate lock token: %w", err)
	}

	// Lock can be safely acquired. Set new expiry time with the token.
	a.locks[key] = inMemoryRecord{
		expiry: now.Add(expiration),
		token:  token,
	}
	return token, nil
}

// ReleaseLock releases an in-memory lock if the provided token matches the current owner.
func (a *InMemoryLockAdapter) ReleaseLock(ctx context.Context, key, token string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	rec, exists := a.locks[key]
	if !exists {
		return errors.New("no lock found to release")
	}
	// Check ownership
	if rec.token != token {
		return errors.New("cannot release lock: token mismatch")
	}
	// Release
	delete(a.locks, key)
	return nil
}
