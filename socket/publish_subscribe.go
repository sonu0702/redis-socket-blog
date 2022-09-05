package socket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	storePubsubChannel string = "socketevent"
)

type CacheStore struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

func NewCacheStore(ctx context.Context, client *redis.Client) *CacheStore {
	return &CacheStore{Ctx: ctx, RedisClient: client}
}

func (c *CacheStore) SubscribeSocketChannel(socketClient *Client) {
	ctx := context.Background()
	redisPubsub := c.RedisClient.Subscribe(ctx, storePubsubChannel)
	go func() {
		for msg := range redisPubsub.Channel() {
			switch msg.Channel {
			case storePubsubChannel:
				fmt.Println("received pubsub message:", msg.Payload)
				publishToSocketSubscribers(socketClient, msg.Payload)
			}
		}
	}()
}

type StorePublishPayload struct {
	AppId       string      `json:"appId"`
	WorkspaceId string      `json:"workspaceId"`
	Payload     interface{} `json:"payload"`
}

type StorePublshedPayload struct {
	Event   ApplicationEventType `json:"event"`
	Payload interface{}          `json:"payload"`
	UserId  string               `json:"userId"`
}
type StorePublishedData struct {
	AppId       string               `json:"appId"`
	WorkspaceId string               `json:"workspaceId"`
	Payload     StorePublshedPayload `json:"payload"`
}

func UnMarshalledPublishedData(publishedMessage string) (*StorePublishedData, error) {
	var newStorePublishPayload StorePublishedData
	err := json.Unmarshal([]byte(publishedMessage), &newStorePublishPayload)
	if err != nil {
		return nil, err
	}
	return &newStorePublishPayload, nil
}

func publishToSocketSubscribers(socketClient *Client, publishedMessage string) error {
	publishedPayload, err := UnMarshalledPublishedData(publishedMessage)
	if err != nil {
		return err
	}
	fmt.Println("Broadcast event : ", publishedPayload.Payload.Event, " to workspace : ", publishedPayload.WorkspaceId, "app :", publishedPayload.AppId)
	socketClient.BroadCast(publishedPayload.WorkspaceId, publishedPayload.AppId,
		publishedPayload.Payload)
	return nil
}

func getPublishPayload(workspaceId, appId string, payload interface{}) (string, error) {
	newStorePublishPayload := StorePublishPayload{
		WorkspaceId: workspaceId,
		AppId:       appId,
		Payload:     payload,
	}
	publishPayloadByte, err := json.Marshal(newStorePublishPayload)
	if err != nil {
		return "", err
	}
	return string(publishPayloadByte), nil
}

func (c *CacheStore) PublishToSocketChannel(message string) error {
	redisCommand := c.RedisClient.Publish(c.Ctx, storePubsubChannel, message)
	_, err := redisCommand.Result()
	if err != nil {
		return err
	}
	return nil
}

func Publish(ctx context.Context, redisClient *redis.Client,
	workspaceId, appId string, payload interface{}) error {
	publishableData, err := getPublishPayload(workspaceId, appId, payload)
	if err != nil {
		return err
	}
	redisStore := NewCacheStore(ctx, redisClient)
	err = redisStore.PublishToSocketChannel(publishableData)
	if err != nil {
		return err
	}
	return nil
}
