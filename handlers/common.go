package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type ReqQueryParameters struct {
	API string
	ID  string
}

const (
	APIStations = "stations"
	APIJourneys = "journeys"
)

type APIHandler struct {
	lastRequestTime time.Time
	newRequest      ReqQueryParameters
}

var handlerMap = map[string]func(ReqQueryParameters) (string, error){
	APIStations: StationsGet,
	APIJourneys: JourneysGet,
}

func (h *APIHandler) HandleRequest(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	if time.Since(h.lastRequestTime) < time.Second {
		return createErrorResponse(http.StatusTooManyRequests, "please wait 1 second before making another request"), nil
	}

	h.lastRequestTime = time.Now()

	queryParameters := request.QueryStringParameters

	if id, ok := queryParameters["id"]; ok && id != "" {
		if id == "0" {
			return createErrorResponse(http.StatusBadRequest, "0 is an invalid ID"), nil
		}

		h.newRequest.ID = id
	}

	api, ok := queryParameters["api"]
	if !ok || api == "" {
		return createErrorResponse(http.StatusBadRequest, "missing required query parameter: api"), nil
	}

	if queryParameters["id"] == "" {
		h.newRequest.ID = "0"
	}
	h.newRequest.API = api

	result, resultError := handleAPIRequest(h.newRequest.API, h.newRequest)
	if resultError != nil {
		return createErrorResponse(http.StatusBadRequest, string([]byte(resultError.Error()))), nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(result)
	return response, nil
}

func handleAPIRequest(api string, request ReqQueryParameters) (string, error) {
	handler, ok := handlerMap[api]
	if !ok {
		return "", errors.New(api + " api does not exist, try stations or journeys.")
	}
	return handler(request)
}

func createErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       message,
	}
}
