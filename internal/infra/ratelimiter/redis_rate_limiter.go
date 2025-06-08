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
	"log"
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

func (r *RedisStorage) CheckAndIncrement(ctx context.Context, key string, maxRequests int64, lockout time.Duration) (bool, int, error) {
	// 1. Check if the key is locked out
	locked, err := r.IsLockedOut(ctx, key)
	if err != nil {
		return false, 0, err
	}
	if locked {
		return false, 0, nil
	}

	requestKey := requestsPrefix + key

	// 2. Increment the request count
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, requestKey)
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, 0, err
	}

	count := incr.Val()

	// 3. Set expiration only on the first hit (1-second window)
	if count == 1 {
		r.client.Expire(ctx, requestKey, time.Second)
	}

	// 4. Determine if the request is allowed
	allowed := maxRequests == 0 || count <= maxRequests

	// 5. Optionally: apply lockout if limit is exceeded
	if !allowed {
		r.SetLockOut(ctx, key, lockout)
	}

	log.Println("CheckAndIncrement - key:", key, "count:", count, "maxRequests:", maxRequests, "allowed:", allowed)
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

func (r *RedisStorage) Reset(ctx context.Context, key string) error {
	pipe := r.client.Pipeline()
	pipe.Del(ctx, requestsPrefix+key)
	pipe.Del(ctx, lockoutPrefix+key)
	_, err := pipe.Exec(ctx)
	return err
}
