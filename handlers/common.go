package handlers

import (
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type ReqQueryParameters struct {
	API    string
	ID     string
	Method string
}

const (
	APIStations = "stations"
	APIJourneys = "journeys"
)

type APIHandler struct {
	lastRequestTime time.Time
	newRequest      ReqQueryParameters
}

var handlerMap = map[string]func(ReqQueryParameters) events.APIGatewayProxyResponse{
	APIStations: func(request ReqQueryParameters) events.APIGatewayProxyResponse {
		if request.Method == "GET" {
			return StationsGet(request)
		} else if request.Method == "POST" {
			return StationsPost(request)
		} else {
			return createErrorResponse(http.StatusMethodNotAllowed, "only GET and POST requests are allowed for stations api")
		}
	},
	APIJourneys: func(request ReqQueryParameters) events.APIGatewayProxyResponse {
		return JourneysGet(request)
	},
}

func (h *APIHandler) HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	if time.Since(h.lastRequestTime) < time.Second {
		return createErrorResponse(http.StatusTooManyRequests, "please wait 1 second before making another request"), nil
	}

	h.lastRequestTime = time.Now()

	queryParameters := request.QueryStringParameters

	api, ok := queryParameters["api"]
	if !ok || api == "" {
		return createErrorResponse(http.StatusBadRequest, "missing required query parameter: api"), nil
	}

	h.newRequest.API = api

	if id, ok := queryParameters["id"]; ok && id != "" {
		if id == "0" {
			return createErrorResponse(http.StatusBadRequest, "0 is an invalid ID"), nil
		}

		h.newRequest.ID = id
	}

	if queryParameters["id"] == "" {
		h.newRequest.ID = "0"
	}

	h.newRequest.Method = request.HTTPMethod

	return handleAPIRequest(h.newRequest.API, h.newRequest), nil
}

func handleAPIRequest(api string, request ReqQueryParameters) events.APIGatewayProxyResponse {
	handler, ok := handlerMap[api]
	if !ok {
		return createErrorResponse(http.StatusBadRequest, api+" api does not exist, try stations or journeys.")
	}

	return handler(request)
}

func createErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
	}
}
