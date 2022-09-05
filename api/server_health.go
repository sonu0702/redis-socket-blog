package api

import (
	"encoding/json"
	"net/http"
)

//HealthCheckResponse health check
type HealthCheckResponse struct {
	Health string `json:"health,omitempty"`
}

//WebappLogSink logs from webapp
type WebappLogSink struct {
	LogType string                 `json:"logType,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

func (client *Client) health(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Health: "ok",
	}
	json.NewEncoder(w).Encode(*NewResponse(SUCCESS, response, nil))
}
