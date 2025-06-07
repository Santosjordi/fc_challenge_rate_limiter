// Package db provides a Redis-backed implementation for rate limiting functionality.
// It defines the RedisStorage struct, which interacts with a Redis server to manage
// request counts, lockout states, and rate limit configurations for individual keys.
//
// The RedisStorage struct offers methods to:
//   - Check the health of the Redis connection (Ping)
//   - Increment and check request counts with respect to configured rate limits (CheckAndIncrement)
//   - Determine and set lockout status for keys (IsLockedOut, SetLockOut)
//   - Retrieve current rate limit status and configuration (GetRateLimit)
//   - Reset request and lockout states for a key (Reset)
//
// The package uses Redis pipelines for efficient multi-command execution and stores
// configuration and state using key prefixes for requests, lockouts, and configs.
// Default rate limiting values are provided if no configuration is found for a key.
package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	requestsPrefix = "requests:"
	lockoutPrefix  = "lockout:"
	configPrefix   = "config:"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

// Ping checks the connection to the Redis server.
// Returns nil if the connection is healthy, otherwise returns an error.
func (r *RedisStorage) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

// CheckAndIncrement checks if the number of requests associated with the given key
// has exceeded the allowed limit within a time window. It increments the request
// count for the key in Redis and retrieves the maximum allowed requests. If this is
// the first request, it sets an expiration for the key. The function returns whether
// the request is allowed, the number of remaining requests, and any error encountered.
//
// Parameters:
//
//	ctx - context for controlling cancellation and deadlines.
//	key - unique identifier for the rate limit.
//
// Returns:
//
//	allowed - true if the request is within the allowed limit, false otherwise.
//	remaining - number of requests remaining before reaching the limit.
//	err - error encountered during the operation, if any.
func (r *RedisStorage) CheckAndIncrement(ctx context.Context, key string) (bool, int, error) {
	requestKey := requestsPrefix + key

	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, requestKey)
	limit := pipe.Get(ctx, configPrefix+key)
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, 0, err
	}

	count := incr.Val()
	maxRequests := int64(0)
	if limit.Err() == nil {
		maxRequests, _ = limit.Int64()
	}

	// If this is first request, set expiration
	if count == 1 {
		r.client.Expire(ctx, requestKey, time.Minute)
	}

	allowed := maxRequests == 0 || count <= maxRequests
	return allowed, int(maxRequests - count), nil
}

func (r *RedisStorage) IsLockedOut(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, lockoutPrefix+key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisStorage) SetLockOut(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, lockoutPrefix+key, true, duration).Err()
}

func (r *RedisStorage) GetRateLimit(ctx context.Context, key string) (RateLimit, error) {
	var rateLimit RateLimit

	// Get current request count
	requestCount, err := r.client.Get(ctx, requestsPrefix+key).Int()
	if err != nil && err != redis.Nil {
		return rateLimit, err
	}

	// Get lockout status
	isLocked, ttl, err := r.getLockoutStatus(ctx, key)
	if err != nil {
		return rateLimit, err
	}

	// Get stored config
	configJSON, err := r.client.Get(ctx, configPrefix+key).Result()
	if err != nil && err != redis.Nil {
		return rateLimit, err
	}

	if err == redis.Nil {
		// Return default config
		return RateLimit{
			WindowSize:      time.Minute,
			MaxRequests:     60,
			LockoutDuration: time.Minute * 5,
			CurrentRequests: requestCount,
			IsLocked:        isLocked,
			LockedUntil:     time.Now().Add(ttl),
		}, nil
	}

	if err := json.Unmarshal([]byte(configJSON), &rateLimit); err != nil {
		return rateLimit, err
	}

	rateLimit.CurrentRequests = requestCount
	rateLimit.IsLocked = isLocked
	rateLimit.LockedUntil = time.Now().Add(ttl)

	return rateLimit, nil
}

func (r *RedisStorage) Reset(ctx context.Context, key string) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, requestsPrefix+key)
	pipe.Del(ctx, lockoutPrefix+key)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisStorage) getLockoutStatus(ctx context.Context, key string) (bool, time.Duration, error) {
	pipe := r.client.Pipeline()
	exists := pipe.Exists(ctx, lockoutPrefix+key)
	ttl := pipe.TTL(ctx, lockoutPrefix+key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}
	return exists.Val() > 0, ttl.Val(), nil
}
