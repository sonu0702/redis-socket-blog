package server

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func connectRedis(ctx context.Context, options *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(options)
	_, err := client.Ping(ctx).Result()
	return client, err
}
