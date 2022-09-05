package socket

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type ApplicationEventType string

const (
	Subscribe   ApplicationEventType = "SUBSCRIBE"
	UnSubscribe ApplicationEventType = "UNSUBSCRIBE"
	//app
	UpdateApp ApplicationEventType = "UPDATE_APP"
)

const (
	roomNameFormat string = "workspace_%s_app_%s"
)

type ApplicationEventPayload struct {
	WorkspaceId string `json:"workspaceId"`
	AppId       string `json:"appId"`
	Token       string `json:"token"`
}

type ApplicationSubscriptionPayload struct {
	Event   ApplicationEventType `json:"event"`
	Payload ApplicationEventPayload  `json:"payload"`
}

type ApplicationSubscriptionResponsePayload struct {
	Event   ApplicationEventType `json:"event"`
	Payload ApplicationEventPayload  `json:"payload"`
	Status  string               `json:"status"`
}

func NewApplicationSubscriptionResponse(status string, event ApplicationEventType,
	payload ApplicationEventPayload) ApplicationSubscriptionResponsePayload {
	return ApplicationSubscriptionResponsePayload{Status: status,
		Event: event, Payload: payload}
}

// start app
type AppUpdatePayload struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	WorkspaceID string `json:"workspaceID,omitempty"`
}
type SocketUpdateApp struct {
	Event   ApplicationEventType `json:"event"`
	Payload AppUpdatePayload     `json:"payload"`
	UserId  string               `json:"userId"`
}

func NewUpdateAppPayload(userId string, appID string, app interface{}) (*SocketUpdateApp, error) {
	var newApp AppUpdatePayload
	appByte, err := json.Marshal(app)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(appByte, &newApp)
	if err != nil {
		return nil, err
	}
	newApp.ID = appID
	return &SocketUpdateApp{Event: UpdateApp, Payload: newApp, UserId: userId}, nil
}

type PayloadOptions struct {
	EventType   ApplicationEventType
	WorkspaceId string
	AppId       string
	UserId      string
	Id          string
	Data        interface{}
	OnlineUsers []OnlineUserMetadata
}

func Notify(ctx context.Context, redisClient *redis.Client, payloadOptions PayloadOptions) {
	switch payloadOptions.EventType {
	case UpdateApp:
		socketUpdatePayload, _ := NewUpdateAppPayload(payloadOptions.UserId, payloadOptions.AppId,
			payloadOptions.Data)
		if socketUpdatePayload != nil {
			Publish(ctx, redisClient, payloadOptions.WorkspaceId, payloadOptions.AppId, socketUpdatePayload)
		}
	}
}
