package socket

import (
	"context"
	"time"
)

//Check the request time
func (c *Client) newContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	return ctx, cancel
}
