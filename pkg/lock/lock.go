package lock

import (
	"context"
	"fmt"
	"time"
)

/*
Package distributedlock provides a robust and easy-to-use distributed locking mechanism. It is designed
to prevent concurrency issues (e.g., race conditions, deadlocks) by allowing only one holder for a lock
at a time across distributed processes or services.

The package follows an adapter-port pattern:
- LockClient is the "port" that the application uses.
- LockAdapter is the interface defining "adapter" implementations (e.g., Redis, in-memory).
*/

// Adapter defines methods that any distributed lock implementation must provide.
// Each implementation is responsible for managing lock acquisition, release,
// and ensuring expiry or refresh policies as needed.
type Adapter interface {
	// AcquireLock attempts to acquire a lock on the given key with the specified expiration.
	// If the lock is already held by another process, AcquireLock may either block, retry,
	// or fail immediately depending on the adapter's logic and the context's deadline.
	//
	// Returns a unique "token" (string) that identifies the holder of the lock. The caller
	// must provide this token when releasing the lock to ensure only the owner can release it.
	AcquireLock(ctx context.Context, key string, expiration time.Duration) (string, error)

	// ReleaseLock releases the lock for the given key if the provided token matches
	// the token of the current lock owner. If the lock is held by a different token,
	// it should fail.
	ReleaseLock(ctx context.Context, key, token string) error
}

// Locker is the primary interface for end users. It leverages a LockAdapter
// to perform the actual lock operations.
type Locker struct {
	adapter Adapter
}

// NewLocker creates a new Locker with the provided adapter.
func NewLocker(adapter Adapter) Locker {
	return Locker{
		adapter: adapter,
	}
}

// Acquire requests ownership of the lock for the specified key. If another client
// already holds the lock, Acquire will fail (or potentially block or retry depending
// on the adapter implementation). The lock will expire automatically after the
// provided duration if not refreshed. Returns a token that uniquely identifies
// the lock holder; you must keep this token to release the lock later.
func (lc Locker) Acquire(ctx context.Context, key string, expiration time.Duration) (string, error) {
	token, err := lc.adapter.AcquireLock(ctx, key, expiration)
	if err != nil {
		return "", fmt.Errorf("failed to acquire lock for key=%q: %w", key, err)
	}
	return token, nil
}

// Release relinquishes ownership of the lock for the specified key if the provided
// token matches the lock's current owner. If the token does not match, release fails.
func (lc Locker) Release(ctx context.Context, key, token string) error {
	if err := lc.adapter.ReleaseLock(ctx, key, token); err != nil {
		return fmt.Errorf("failed to release lock for key=%q: %w", key, err)
	}
	return nil
}
