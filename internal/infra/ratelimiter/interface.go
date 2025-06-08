package db

import (
	"context"
	"time"
)

type RateLimiter interface {
	// CheckAndIncrement atomically checks and increments counter
	// Returns: allowed bool, remaining int, error
	CheckAndIncrement(ctx context.Context, key string, maxRequests int64, lockout time.Duration) (bool, int, error)

	// IsLockedOut checks if key is in lockout state
	IsLockedOut(ctx context.Context, key string) (bool, error)

	// SetLockOut puts key in lockout state with expiration
	SetLockOut(ctx context.Context, key string, duration time.Duration) error

	// Reset resets counters for a key
	Reset(ctx context.Context, key string) error
}
