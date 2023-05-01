package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	API    string
	Method string
}

type ReqQueryParameters struct {
	API string
	ID  string
}

var handlersMap = map[string]Handler{
	"stations": {
		API:    "stations",
		Method: "GET",
	},
	"journeys": {
		API:    "journeys",
		Method: "GET",
	},
}

type APIHandler struct {
	lastRequestTime time.Time
	newRequest      ReqQueryParameters
}

func (h *APIHandler) HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {

	// Check if enough time has passed since the last request
	if time.Since(h.lastRequestTime) < time.Second {
		response.StatusCode = http.StatusTooManyRequests
		response.Body = "please wait 1 second before making another request"
		return
	}

	// Update last request time to now
	h.lastRequestTime = time.Now()

	queryParameters := request.QueryStringParameters

	if id, ok := queryParameters["id"]; ok && id != "" {
		if id == "0" {
			response.StatusCode = http.StatusBadRequest
			response.Body = "0 is an invalid ID"
			fmt.Println("response", response)
			return
		}
		h.newRequest.ID = id
	}

	api, ok := queryParameters["api"]
	if !ok || api == "" {
		response.StatusCode = http.StatusBadRequest
		response.Body = "missing required query parameter: api"
		return
	}

	if queryParameters["id"] == "" {
		h.newRequest.ID = "0"
	}
	h.newRequest.API = api

	handle, ok := handlersMap[h.newRequest.API]
	if !ok || handle.Method != "GET" {
		response.StatusCode = http.StatusMethodNotAllowed
		response.Body = "method not allowed"
		return
	}

	result, stationsError := StationsGet(h.newRequest)
	if stationsError != nil {
		response.Body = string([]byte(stationsError.Error()))
		response.StatusCode = 500
		return
	}

	response.StatusCode = http.StatusOK
	response.Body = string(result)
	return
}
