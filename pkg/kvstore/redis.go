package kvstore

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisStore{Client: rdb}
}

func (r *RedisStore) Ping() error {
	ctx := context.Background()
	return r.Client.Ping(ctx).Err()
}

func (r *RedisStore) Set(key string, value any) error {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, 0).Err()
}

func (r *RedisStore) Get(key string) (any, error) {
	ctx := context.Background()
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisStore) Exist(key string) (bool, error) {
	_, err := r.Get(key)
	return err == nil, nil
}

func (r *RedisStore) Delete(key string) error {
	ctx := context.Background()
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisStore) Close() error {
	return r.Client.Close()
}
