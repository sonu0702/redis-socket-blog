package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	health  = "/health"
	sockets = "/socket.io"
	apps    = "/apps"
)

func (client *Client) registerOpenAppRoutes(router *mux.Router) {
	usersAPI := router.PathPrefix(apps).Subrouter()
	usersAPI.HandleFunc("/", client.updateApp).Methods(http.MethodPut)
}

func (client *Client) registerOpenServerHealthCheck(router *mux.Router) {
	healthAPI := router.PathPrefix(health).Subrouter()
	healthAPI.HandleFunc("/", client.health).Methods(http.MethodGet)
}

func (client *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	fmt.Printf("client.SocketClient %+v", client.SocketClient.SocketServer)
	client.SocketClient.SocketServer.ServeHTTP(w, r)
}

func (client *Client) registerSockets(router *mux.Router) {
	socketsAPI := router.PathPrefix(sockets).Subrouter()
	socketsAPI.HandleFunc("/", client.ServeHTTP).Methods(http.MethodGet, http.MethodPost)
}

func (client *Client) registerOpenRoutes(router *mux.Router) {
	router.Use(client.setResponseHeader)
	router.Use(client.panicRecovery)
	client.registerOpenServerHealthCheck(router)
	client.registerOpenAppRoutes(router)
}

func (client *Client) registerSocketRoutes(router *mux.Router) {
	client.registerSockets(router)
}
