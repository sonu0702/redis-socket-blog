package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	appError "redis-blog/error"
	"redis-blog/socket"
)

type UpdateAppRequest struct {
	WorkspaceId string `json:"workspaceId,omitempty"`
	AppId       string `json:"appId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

func (client *Client) updateApp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newRequestContext(r)
	defer cancel()
	var appRequest UpdateAppRequest
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		setBadRequestErrorResponse(w, appError.GenerateErrorCode(appError.App,
			appError.InvalidPayloadErrorCode))
	}
	defer r.Body.Close()
	err = json.Unmarshal(requestBody, &appRequest)
	if err != nil {
		setBadRequestErrorResponse(w, appError.GenerateErrorCode(appError.App,
			appError.UnMarshallingErrorCode))
		return
	}

	response := map[string]interface{}{
		"displayName": appRequest.DisplayName,
		"workspaceId" :appRequest.WorkspaceId,
		"appId":appRequest.AppId,
	}
	//notify
	payloadOptions := socket.PayloadOptions{EventType: socket.UpdateApp, WorkspaceId: appRequest.WorkspaceId,
		AppId: appRequest.AppId, Data: appRequest}
	socket.Notify(ctx, client.Redis, payloadOptions)
	json.NewEncoder(w).Encode(*NewResponse(SUCCESS, response, nil))
}
