package api

import (
	"context"
	"net/http"
	"time"
)

//Check the request time
func newRequestContext(r *http.Request) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Second)
	return ctx, cancel
}
