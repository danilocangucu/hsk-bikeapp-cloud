package handlers

import (
	"encoding/json"
	"fmt"
	db "hsk-bikeapp-solita-cloud/database"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

var DB db.Db

func StationsGet(newRequest ReqQueryParameters) events.APIGatewayProxyResponse {
	var err error

	if newRequest.ID < "0" {
		log.Println("invalid query parameter:", newRequest.ID)
		return createErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v is an invalid ID", newRequest.ID))
	}

	DB, err = db.OpenDatabase()
	if err != nil {
		log.Println("Error opening database:", err)
		return createErrorResponse(http.StatusInternalServerError, "an error has been produced trying to access the database")
	}
	defer DB.CloseDatabase()

	filter := db.StationFilter{}
	if newRequest.ID != "" {
		filter.ID, err = strconv.Atoi(newRequest.ID)
		if err != nil {
			log.Println("Invalid query parameter:", err)
			return createErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v is an invalid ID", err))
		}
	}

	var station db.Station
	var stations []db.Station
	if filter.ID != 0 {
		station, err = DB.GetSingleStation(filter)
		if err != nil {
			log.Printf("Error while getting station ID %v: %v", filter.ID, err)
			return createErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error while getting station ID %v", filter.ID))
		}
	} else if filter.ID == 0 {
		stations, err = DB.GetAllStations()
		if err != nil {
			log.Println("Error while getting stations:", err)
			return createErrorResponse(http.StatusInternalServerError, "error while getting stations")
		}
	}

	var responseJSON []byte
	if station != (db.Station{}) {
		responseJSON, err = json.Marshal(station)
		if err != nil {
			log.Println("error while marshalling station:", err)
			return createErrorResponse(http.StatusInternalServerError, "oops! something went wrong while processing your station request")
		}
	} else {
		responseJSON, err = json.Marshal(stations)
		if err != nil {
			log.Println("error while marshalling stations:", err)
			return createErrorResponse(http.StatusInternalServerError, "oops! something went wrong while processing your stations request")
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}
}

func StationsPost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error

	DB, err = db.OpenDatabase()
	if err != nil {
		log.Println("Error opening database:", err)
		return createErrorResponse(http.StatusInternalServerError, "an error has been produced trying to access the database")
	}
	defer DB.CloseDatabase()

	var newStation db.Station
	if err = json.Unmarshal([]byte(request.Body), &newStation); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("invalid request body: %v", err),
		}
	}

	validationErrors := DB.ValidateNewStation(newStation)
	if len(validationErrors) > 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("invalid new station: %v", validationErrors),
		}
	}

	if err = DB.AddNewStation(newStation); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("failed to add new station: %v", err),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "New station added successfully!",
	}
}
