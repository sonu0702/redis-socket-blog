package server

import (
	"context"
	"redis-blog/socket"
	"net/http"

	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func initSocketServer() *socketio.Server {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	return server
}

func initSocketChannelSubscriber(ctx context.Context, redisClient *redis.Client,
	socketClient *socket.Client) {
	newStore := socket.NewCacheStore(ctx, redisClient)
	newStore.SubscribeSocketChannel(socketClient)
}
