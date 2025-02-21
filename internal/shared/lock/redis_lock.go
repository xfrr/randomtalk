package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// ErrLockAlreadyHeld is returned when a lock is already held.
var ErrLockAlreadyHeld = errors.New("lock already held")

// ErrLockAlreadyReleased is returned when a lock is already released.
var ErrLockAlreadyReleased = errors.New("lock already released")

// =============================================================================
// Redis Adapter
// =============================================================================

// RedisClient defines the methods that a Redis client must implement.
type RedisClient interface {
	// SetNX sets key to hold the string value if key does not exist.
	// Returns true if the key was set, false otherwise.
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	// Get retrieves the value of a key.
	Get(ctx context.Context, key string) (string, error)

	// Del deletes one or more keys.
	// Returns the number of keys removed.
	Del(ctx context.Context, keys ...string) (int64, error)

	// Eval is used to evaluate a Lua script in Redis. You can use this to ensure an
	// atomic compare-and-del if desired. For demonstration, we're leaving it commented.
	// Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}

type RedisLockAdapter struct {
	client *redis.Client
}

// NewRedisLockAdapter creates a new RedisLockAdapter with the given RedisClient.
func NewRedisLockAdapter(client *redis.Client) *RedisLockAdapter {
	return &RedisLockAdapter{
		client: client,
	}
}

// AcquireLock attempts to atomically set a key with a unique token if it does not exist.
// If the lock is acquired successfully, returns the token. If the lock is already held,
// returns an error.
func (a *RedisLockAdapter) AcquireLock(ctx context.Context, key string, expiration time.Duration) (string, error) {
	// Generate a unique token. We'll store it in Redis if the lock is free.
	token, err := generateRandomToken()
	if err != nil {
		return "", fmt.Errorf("unable to generate token: %w", err)
	}

	bcmd := a.client.SetNX(ctx, key, token, expiration)
	if bcmd.Err() != nil {
		return "", fmt.Errorf("error setting lock key in redis: %w", bcmd.Err())
	}

	if !bcmd.Val() {
		return "", ErrLockAlreadyHeld
	}
	return token, nil
}

// ReleaseLock releases the lock only if the stored token matches the provided token.
// If there's a mismatch, it returns an error. This approach ensures that only the
// owner can release the lock. Below is a simple GET + DEL approach. For a truly
// atomic compare-and-delete, consider using a Lua script via `EVAL`.
func (a *RedisLockAdapter) ReleaseLock(ctx context.Context, key, token string) error {
	strCmd := a.client.Get(ctx, key)
	if strCmd.Err() != nil {
		return fmt.Errorf("error getting lock key in redis: %w", strCmd.Err())
	}

	val, err := strCmd.Result()
	if err != nil {
		return fmt.Errorf("error getting lock key in redis: %w", err)
	}

	// If the lock doesn't exist or tokens mismatch, fail.
	if val == "" {
		return errors.New("no lock found to release")
	}

	if val != token {
		return errors.New("cannot release lock: token mismatch")
	}

	// Now we can safely delete.
	intCmd := a.client.Del(ctx, key)
	if intCmd.Err() != nil {
		return fmt.Errorf("error deleting lock key in redis: %w", intCmd.Err())
	}

	// Check if the key was deleted.
	if intCmd.Val() == 0 {
		return ErrLockAlreadyReleased
	}

	return nil
}
