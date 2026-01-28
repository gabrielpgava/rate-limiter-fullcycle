package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func redisKey(prefix, id string) string {
	return prefix + ":" + id
}

func RedisGet[T any](rdb *redis.Client, prefix, id string) (*T, bool, error) {
	key := redisKey(prefix, id)

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var out T
	if err := json.Unmarshal([]byte(val), &out); err != nil {
		return nil, false, err
	}

	return &out, true, nil
}

func RedisSet[T any](
	rdb *redis.Client,
	prefix, id string,
	value *T,
	ttl time.Duration,
) error {
	key := redisKey(prefix, id)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, data, ttl).Err()
}

func RedisDelete(rdb *redis.Client, prefix, id string) error {
	return rdb.Del(ctx, redisKey(prefix, id)).Err()
}
