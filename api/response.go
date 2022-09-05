package api

import (
	"encoding/json"
	"net/http"
	"os"
	appleError "redis-blog/error"
	"reflect"
)

//Status reponse
type Status string

const (
	//ERROR error status
	ERROR Status = "ERROR"
	//SUCCESS status
	SUCCESS Status = "SUCCESS"
)

//Response api response
type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
	Error  error       `json:"error"`
}

//IsInstanceOf ..
func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

//NewResponse generates API response
func NewResponse(status Status, data interface{}, err error) *Response {
	if err == nil {
		return &Response{Status: status, Data: data, Error: err}
	}
	if !IsInstanceOf(err, &appleError.ApplicationError{}) {
		newResponse := Response{Status: status, Data: data, Error: appleError.New(appleError.InternalServerErrorCode,
			appleError.InternalServerError, err.Error())}
		if os.Getenv("ENVIRONMENT") != "dev" {
			newResponse.Error = appleError.New(appleError.InternalServerErrorCode,
				appleError.InternalServerError, appleError.InternalServerError)
		}
		return &newResponse
	}
	response := Response{Status: status, Data: data, Error: err}
	return &response
}

func setResponse(w http.ResponseWriter, response *Response) {
	json.NewEncoder(w).Encode(response)
}

//NewInternalServerErrorResponse writes server error to response
func setInternalServerErrorResponse(w http.ResponseWriter,
	errorCode string, requestID string) {
	response := NewResponse(ERROR, nil, appleError.New(errorCode,
		appleError.InternalServerError, appleError.InternalServerError))
	setResponse(w, response)
}

//setAppleErrorResponse writes server error to response
func setAppleErrorResponse(w http.ResponseWriter,
	requestID string, err error) {
	response := NewResponse(ERROR, nil, err)
	setResponse(w, response)
}

func setBadRequestErrorResponse(w http.ResponseWriter,
	errorCode string) {
	response := NewResponse(ERROR, nil, appleError.New(errorCode,
		appleError.BadRequestError, appleError.InvalidRequestPayload))
	setResponse(w, response)
}

func (c *Client) setResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		next.ServeHTTP(w, r)
	})
}

//general
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func setCodedBadRequestErrorResponse(w http.ResponseWriter,
	errorCode string, requestID string) {
	response := NewResponse(ERROR, nil, appleError.New(errorCode,
		appleError.BadRequestError, appleError.InvalidRequestPayload))
	w.WriteHeader(http.StatusBadRequest)
	setResponse(w, response)
}

func setCodedInternalServerErrorResponse(w http.ResponseWriter,
	errorCode string, requestID string) {
	response := NewResponse(ERROR, nil, appleError.New(errorCode,
		appleError.BadRequestError, appleError.InvalidRequestPayload))
	w.WriteHeader(http.StatusInternalServerError)
	setResponse(w, response)
}

//setAppleErrorResponse writes server error to response
func setCodedForbiddenErrorResponse(w http.ResponseWriter,
	requestID string, err error) {
	response := NewResponse(ERROR, nil, err)
	w.WriteHeader(http.StatusForbidden)
	setResponse(w, response)
}
