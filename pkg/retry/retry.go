// Package retry provides a simple mechanism for retrying operations with
// configurable attempts and backoff strategies, including optional jitter.
package retry

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Config holds the configuration for retries.
type Config struct {
	MaxAttempts       int              // Maximum number of attempts
	InitialDelay      time.Duration    // Initial delay before the first retry
	MaxDelay          time.Duration    // Upper bound for the delay
	BackoffMultiplier float64          // Multiplier for exponential backoff
	Jitter            bool             // Whether to apply random jitter (±50%)
	IsRetryableError  func(error) bool // Determines if the operation should be retried
}

// Option is a function that applies a configuration option to the Config.
type Option func(*Config)

// WithMaxAttempts sets the maximum number of retry attempts.
func WithMaxAttempts(attempts int) Option {
	return func(c *Config) {
		if attempts < 1 {
			attempts = 1
		}
		c.MaxAttempts = attempts
	}
}

// WithInitialDelay sets the initial delay before the first retry.
func WithInitialDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.InitialDelay = delay
	}
}

// WithMaxDelay sets the maximum delay between retries.
func WithMaxDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.MaxDelay = delay
	}
}

// WithBackoffMultiplier sets the multiplier for exponential backoff.
func WithBackoffMultiplier(multiplier float64) Option {
	return func(c *Config) {
		if multiplier < 1 {
			multiplier = 1
		}
		c.BackoffMultiplier = multiplier
	}
}

// WithJitter enables or disables random jitter (±50% of the delay).
func WithJitter(enable bool) Option {
	return func(c *Config) {
		c.Jitter = enable
	}
}

// WithRetryCondition sets the retry function that determines if the operation should be retried.
func WithRetryCondition(condition func(error) bool) Option {
	return func(c *Config) {
		c.IsRetryableError = condition
	}
}

// defaultConfig defines the default retry configuration.
var defaultConfig = Config{
	MaxAttempts:       3, // Retry up to 3 times by default
	InitialDelay:      100 * time.Millisecond,
	MaxDelay:          5 * time.Second,
	BackoffMultiplier: 2.0,
	Jitter:            true,
	IsRetryableError:  func(err error) bool { return true },
}

// Do runs the given operation with retry logic using a background context.
// It will stop either on success, after exhausting MaxAttempts, or if the
// context is canceled. It returns the last error if all attempts fail.
func Do[T any](operation func() (T, error), opts ...Option) (T, error) {
	return DoContext(context.Background(), operation, opts...)
}

// DoContext is like Do, but accepts a context for cancellation or timeouts.
func DoContext[T any](ctx context.Context, operation func() (T, error), opts ...Option) (T, error) {
	var res T

	// Create a local copy of the default config, then apply all options.
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var err error

	// Start with the initial delay
	currentDelay := cfg.InitialDelay
	if currentDelay <= 0 {
		currentDelay = 1 * time.Millisecond
	}

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		// Check if the context is already canceled
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}

		// Perform the operation
		res, err = operation()
		if err == nil {
			return res, nil // success
		}

		// Check if the error is retryable
		if !cfg.IsRetryableError(err) {
			return res, err
		}

		// If not the last attempt, delay before the next retry
		if attempt < cfg.MaxAttempts {
			delay := currentDelay

			// Optional jitter (±50%)
			if cfg.Jitter {
				jitterFraction := 0.5
				jitterOffset := time.Duration(float64(delay) * jitterFraction)
				jitter := time.Duration(rand.Int63n(int64(jitterOffset)*2) - int64(jitterOffset))
				delay = delay + jitter
				if delay < 0 {
					delay = 0
				}
			}

			// Sleep for the delay or until context is canceled
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return res, ctx.Err()
			}

			// Exponential backoff for the next round
			nextDelay := time.Duration(float64(currentDelay) * cfg.BackoffMultiplier)
			if nextDelay > cfg.MaxDelay {
				nextDelay = cfg.MaxDelay
			}
			currentDelay = nextDelay
		}
	}

	return res, fmt.Errorf("failed after %d attempts: %w", cfg.MaxAttempts, err)
}
