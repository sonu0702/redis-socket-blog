package main

import (
	"context"
	"redis-blog/server"
	"time"
)

// Starts the server.
func main() {
	options := server.LoadOptions()

	initTimeout := 10 * time.Second // 10 seconds
	initCtx, cancel := context.WithTimeout(context.Background(), initTimeout)
	defer cancel()

	server.New(initCtx, options).Start()
}
