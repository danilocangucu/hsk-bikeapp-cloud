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
}

var handlerMap = map[string]func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse{
	APIStations: func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		if request.HTTPMethod == "GET" {
			reqQueryParameters := extractReqQueryParameters(request)
			return StationsGet(reqQueryParameters)
		} else if request.HTTPMethod == "POST" {
			return StationsPost(request)
		} else {
			return createErrorResponse(http.StatusMethodNotAllowed, "only GET and POST requests are allowed for stations api")
		}
	},
	APIJourneys: func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		reqQueryParameters := extractReqQueryParameters(request)
		return JourneysGet(reqQueryParameters)
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

	if id, ok := queryParameters["id"]; ok && id != "" {
		if id == "0" {
			return createErrorResponse(http.StatusBadRequest, "0 is an invalid ID"), nil
		}
	}

	if queryParameters["id"] == "" {
		request.QueryStringParameters["id"] = "0"
	}

	return handleAPIRequest(request), nil
}

func handleAPIRequest(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	api := request.QueryStringParameters["api"]

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

func extractReqQueryParameters(request events.APIGatewayProxyRequest) ReqQueryParameters {
	return ReqQueryParameters{
		API:    request.QueryStringParameters["api"],
		ID:     request.QueryStringParameters["id"],
		Method: request.HTTPMethod,
	}
}
