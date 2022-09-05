package api

import (
	"redis-blog/socket"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)


//Options API Options
type Options struct {
	PathPrefix            string
	SocketRoutePathPrefix string
	RedisClient           *redis.Client
	SocketClient          *socket.Client
}

//API Router
type API struct {
	Router *mux.Router
}

//Client contains all the client objects
type Client struct {
	Redis        *redis.Client
	SocketClient *socket.Client
}

//New initializes the API.
func New(options *Options) *API {
	router := mux.NewRouter()
	apiV1 := router.PathPrefix(options.PathPrefix).Subrouter()
	apiSocketV1 := router.PathPrefix(options.SocketRoutePathPrefix).Subrouter()
	api := &API{Router: router}
	client := Client{
		Redis:        options.RedisClient,
		SocketClient: options.SocketClient,
	}

	client.registerOpenRoutes(apiV1)
	client.registerSocketRoutes(apiSocketV1)
	return api
}
