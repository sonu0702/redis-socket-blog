package socket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

func (c *CacheStore) set(key string, value interface{}, expiration time.Duration) error {
	statusCmd := c.RedisClient.Set(c.Ctx, key, value, expiration)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func (c *CacheStore) get(key string) (string, error) {
	redisCmd := c.RedisClient.Get(c.Ctx, key)
	value, err := redisCmd.Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func getRoomUsers(ctx context.Context, redisClient *redis.Client, roomName string) (map[string]OnlineUserMetadata, error) {
	redisStore := NewCacheStore(ctx, redisClient)
	existingUsersString, err := redisStore.get(roomName)
	if err != nil && err.Error() == string(redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var existingUsers map[string]OnlineUserMetadata
	err = json.Unmarshal([]byte(existingUsersString), &existingUsers)
	if err != nil {
		return nil, err
	}
	return existingUsers, nil
}

func setRoomUsers(ctx context.Context, redisClient *redis.Client, roomName string,
	roomUsers map[string]OnlineUserMetadata) error {
	roomUsersByte, err := json.Marshal(roomUsers)
	if err != nil {
		return err
	}
	redisStore := NewCacheStore(ctx, redisClient)
	err = redisStore.set(roomName, string(roomUsersByte), 0)
	if err != nil {
		return err
	}
	return nil
}
