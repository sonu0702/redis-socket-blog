package socket

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
)

type Client struct {
	SocketServer *socketio.Server
	RedisClient  *redis.Client
}

type OnlineUserMetadata struct {
	UserId      string `json:"userId"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	PhotoURL    string `json:"photoURL"`
	Color       string `json:"color"`
	SessionId   string `json:"sessionId"`
}
type SessionContext struct {
	WorkspaceId string
	AppId       string
	OnlineUser  OnlineUserMetadata
	SessionId   string
}

func New(redisClient *redis.Client,
	socketServer *socketio.Server) *Client {
	return &Client{RedisClient: redisClient,
		SocketServer: socketServer}
}

func colorPicker() string {
	colors := []string{"#FF6F2C", "#20C933", "#FFDCE5", "#FFEAB6", "#EDE2FE"}
	value := rand.Intn(len(colors))
	return colors[value]
}

func (c *Client) authenticate(next func(socketio.Conn, ApplicationSubscriptionPayload) error) func(socketio.Conn, ApplicationSubscriptionPayload) error {
	return func(conn socketio.Conn, msg ApplicationSubscriptionPayload) error {
		fmt.Printf("SocketAuthenticating====================================== %s\n", conn.ID())

		newSessionContext := SessionContext{
			SessionId:   conn.ID(),
			WorkspaceId: msg.Payload.WorkspaceId,
			AppId:       msg.Payload.AppId,
			OnlineUser: OnlineUserMetadata{
				UserId:      "",
				Email:       "",
				PhotoURL:    "",
				DisplayName: "",
				Color:       colorPicker(),
				SessionId:   conn.ID(),
			},
		}
		fmt.Printf("Subscriber  is Authenticated\n")
		conn.SetContext(newSessionContext)
		return next(conn, msg)
	}
}

func (c *Client) authorize(next func(socketio.Conn, ApplicationSubscriptionPayload) error) func(socketio.Conn, ApplicationSubscriptionPayload) error {
	return func(conn socketio.Conn, msg ApplicationSubscriptionPayload) error {
		sessionContext := conn.Context().(SessionContext)
		fmt.Println("Authorizing...", msg.Payload.WorkspaceId, msg.Payload.AppId, strings.TrimSpace(sessionContext.OnlineUser.UserId))
		fmt.Printf("Subscriber :  %s (%s) is Authorized\n", sessionContext.OnlineUser.UserId,
			sessionContext.OnlineUser.DisplayName)
		return next(conn, msg)
	}
}

func (c *Client) InitEvents() {
	c.SocketServer.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected:", s.ID())
		return nil
	})
	c.SocketServer.OnEvent("/", "application",
		c.authenticate(c.authorize(func(s socketio.Conn, payload ApplicationSubscriptionPayload) error {
			log.Println("roomName : ", getRoomName(payload.Payload.WorkspaceId, payload.Payload.AppId),
				"payload : ", payload.Event)
			roomName := getRoomName(payload.Payload.WorkspaceId, payload.Payload.AppId)
			if payload.Event == Subscribe {
				sessionContext := getSessionContext(s)
				s.Join(roomName)
				log.Printf("subscribed : (%s), sessionId : %s, roomName : %s , allRooms : %v ",
					sessionContext, s.ID(),
					roomName, s.Rooms())
				newApplicationSubscriptionResponse := NewApplicationSubscriptionResponse("SUCCESS",
					payload.Event, payload.Payload)
				s.Emit("application_response", newApplicationSubscriptionResponse)
				return nil
			}
			log.Printf("%s un subscribed to roomName : %s , allRooms : %v ", s.ID(),
				roomName, s.Rooms())
			s.Leave(getRoomName(payload.Payload.WorkspaceId, payload.Payload.AppId))
			newNodezapSubscriptionResponse := NewApplicationSubscriptionResponse("SUCCESS",
				payload.Event, payload.Payload)
			s.Emit("application_response", newNodezapSubscriptionResponse)
			return nil
		})))
	c.SocketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Disconnected ", s.ID(), "reason", reason)
		s.Leave(reason)
		fmt.Println("Disconnected ", s.ID())
	})
}

func getSessionContext(s socketio.Conn) *SessionContext {
	newContext := s.Context()
	if newContext == nil {
		return nil
	}
	if _, isString := newContext.(string); isString {
		return nil
	}
	newSessionContext := newContext.(SessionContext)
	return &newSessionContext
}

func (c *Client) BroadCast(workspaceID, appID string, data interface{}) {
	//broadcast to room
	roomName := getRoomName(workspaceID, appID)
	fmt.Println("publish to :", roomName)
	c.SocketServer.BroadcastToRoom("/", roomName, roomName, data)
}

func getRoomName(workspaceId, appId string) string {
	return fmt.Sprintf(roomNameFormat, workspaceId, appId)
}
