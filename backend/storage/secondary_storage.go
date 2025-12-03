package storage

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type CustomRedisSecondaryStorage struct {
	client *redis.Client
}

func NewCustomRedisSecondaryStorage() *CustomRedisSecondaryStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	return &CustomRedisSecondaryStorage{
		client: rdb,
	}
}

// Get retrieves the value for the given key.
func (s *CustomRedisSecondaryStorage) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Key does not exist
		}
		return nil, err
	}
	return val, nil
}

// Set sets the value for the given key with optional TTL (in seconds).
func (s *CustomRedisSecondaryStorage) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	var expiration time.Duration
	if ttlSeconds > 0 {
		expiration = time.Duration(ttlSeconds) * time.Second
	} else {
		expiration = 0 // No expiration
	}

	return s.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes the value for the given key.
func (s *CustomRedisSecondaryStorage) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
