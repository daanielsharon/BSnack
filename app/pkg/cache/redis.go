package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func DeleteRedisKeysByPattern(ctx context.Context, redisClient *redis.Client, pattern string) error {
	keys, err := redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("redis KEYS failed: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	if _, err := redisClient.Del(ctx, keys...).Result(); err != nil {
		return fmt.Errorf("redis DEL failed: %w", err)
	}

	log.Printf("[CACHE] Deleted %d keys with pattern '%s'", len(keys), pattern)
	return nil
}

func SetJSON(ctx context.Context, rdb *redis.Client, key string, value any, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, bytes, ttl).Err()
}

func SetCount(ctx context.Context, rdb *redis.Client, key string, val int64, ttl time.Duration) error {
	return rdb.Set(ctx, key, strconv.FormatInt(val, 10), ttl).Err()
}
