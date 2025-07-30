package db

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Address string
}

func NewRedisClient(cfg RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}
