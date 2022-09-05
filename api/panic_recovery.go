package api

import (
	"encoding/json"
	"net/http"
	appError "redis-blog/error"
)

func (c *Client) panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				//logging
				newError := appError.New(appError.GlobalPanicErrorCode,
					appError.InternalServerError, appError.InternalServerError)
				json.NewEncoder(w).Encode(*NewResponse(ERROR, nil, newError))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
