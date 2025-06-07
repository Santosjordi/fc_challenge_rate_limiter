package db

import (
	"context"
	"time"
)

type RateLimiter interface {
	// CheckAndIncrement atomically checks and increments counter
	// Returns: allowed bool, remaining int, error
	CheckAndIncrement(ctx context.Context, key string) (bool, int, error)

	// IsLockedOut checks if key is in lockout state
	IsLockedOut(ctx context.Context, key string) (bool, error)

	// SetLockOut puts key in lockout state with expiration
	SetLockOut(ctx context.Context, key string, duration time.Duration) error

	// GetRateLimit gets current limit config and state
	GetRateLimit(ctx context.Context, key string) (RateLimit, error)

	// Reset resets counters for a key
	Reset(ctx context.Context, key string) error
}

type RateLimit struct {
	// Window duration for rate limiting
	WindowSize time.Duration

	// Maximum requests allowed in window
	MaxRequests int

	// Duration to lock out after limit exceeded
	LockoutDuration time.Duration

	// Current state
	CurrentRequests int
	WindowStart     time.Time
	IsLocked        bool
	LockedUntil     time.Time
}
