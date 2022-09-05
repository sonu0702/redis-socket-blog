package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"redis-blog/api"
	"redis-blog/socket"

	// Blank import will autoload env variables from .env file.
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

//APIOptions API options.
type APIOptions struct {
	PathPrefix            string
	SocketRoutePathPrefix string
}

//Options Server options.
type Options struct {
	Redis *redis.Options
	Port  string
	API   APIOptions
}

//Server server object.
type Server struct {
	Options      *Options
	Router       *mux.Router
	RedisClient  *redis.Client
	SocketClient *socket.Client
}

func newCORSoptions() *cors.Cors {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:3000/",
			"http://localhost:8080"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodPatch, http.MethodDelete, http.MethodOptions, http.MethodHead},
		AllowedHeaders: []string{"X-Requested-With", "Content-Type", "Authorization", "Connection", "Host",
			"Upgrade", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Sec-WebSocket-Extensions"},
		AllowCredentials: true,
		Debug:            true,
	})
	return corsOptions
}

//LoadOptions Load server options.
func LoadOptions() *Options {
	var options Options

	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
	}
	options.Redis = &redis.Options{
		Addr:        os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		IdleTimeout: -1,
	}
	if os.Getenv("ENVIRONMENT") != "dev" {
		options.Redis.TLSConfig = &tlsConfig
	}
	options.Port = os.Getenv("PORT")
	if len(options.Port) < 1 {
		log.Fatal("Port not specified.")
	}
	options.API.PathPrefix = "/api/v1"
	options.API.SocketRoutePathPrefix = "/api/socket/v1"
	return &options
}

func connectClients(ctx context.Context, options *Options) *redis.Client {
	var wg errgroup.Group
	var redisClient *redis.Client
	var err error
	wg.Go(func() error {
		redisClient, err = connectRedis(ctx, options.Redis)
		if err != nil {
			fmt.Println("Could not connect redis!")
		}
		return err
	})
	if err := wg.Wait(); err != nil {
		log.Fatal("Error while initializing server: ", err)
	}
	return redisClient
}

//Initialize Initialize.
func initialize(ctx context.Context, options *Options) *Server {
	redisClient := connectClients(ctx, options)
	socketioServer := initSocketServer()
	socketClient := socket.New(redisClient, socketioServer)
	socketClient.InitEvents()
	initSocketChannelSubscriber(ctx, redisClient, socketClient)
	apiOptions := &api.Options{
		PathPrefix:            options.API.PathPrefix,
		SocketRoutePathPrefix: options.API.SocketRoutePathPrefix,
		RedisClient:           redisClient,
		SocketClient:          socketClient,
	}
	api := api.New(apiOptions)
	return &Server{
		Router:       api.Router,
		Options:      options,
		RedisClient:  redisClient,
		SocketClient: socketClient,
	}
}

//New returns an App.
func New(ctx context.Context, options *Options) *Server {
	return initialize(ctx, options)
}

//Start the server.
func (server *Server) Start() {
	go func() {
		if err := server.SocketClient.SocketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.SocketClient.SocketServer.Close()
	http.Handle("/", server.Router)
	log.Println("Server started on port: ", server.Options.Port)
	log.Fatal(http.ListenAndServe(":"+server.Options.Port, newCORSoptions().Handler(server.Router)))
}
