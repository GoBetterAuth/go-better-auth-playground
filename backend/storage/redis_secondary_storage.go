package storage

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSecondaryStorage struct {
	client *redis.Client
}

func NewRedisSecondaryStorage() *RedisSecondaryStorage {
	db := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	return &RedisSecondaryStorage{
		client: db,
	}
}

func (s *RedisSecondaryStorage) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Key does not exist
		}
		return nil, err
	}
	return val, nil
}

func (s *RedisSecondaryStorage) Set(ctx context.Context, key string, value any, ttl *time.Duration) error {
	var expiration time.Duration
	if ttl != nil {
		expiration = time.Duration(*ttl) * time.Second
	} else {
		expiration = 0
	}

	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *RedisSecondaryStorage) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *RedisSecondaryStorage) Incr(ctx context.Context, key string, ttl *time.Duration) (int, error) {
	val, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if ttl != nil && *ttl > 0 {
		_, err := s.client.Expire(ctx, key, *ttl).Result()
		if err != nil {
			return int(val), err
		}
	}

	return int(val), nil
}
