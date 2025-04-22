package db

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

// Register a request from a new token or IP address
func (r *RedisStorage) RegisterRequest(key string) error {
	ctx := context.Background()
	_, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Error incrementing key %s: %v", key, err)
	}
	return nil
}

func (r *RedisStorage) Increment(key string, windowSeconds int) (int, error) {
	ctx := context.Background()
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		r.client.Expire(ctx, key, time.Duration(windowSeconds)*time.Second)
	}

	return int(count), nil
}

func (r *RedisStorage) IsLockedOut(key string) (bool, error) {
	ctx := context.Background()
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (r *RedisStorage) SetLockout(key string, durationSeconds int) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, 1, time.Duration(durationSeconds)*time.Second).Err()
}
