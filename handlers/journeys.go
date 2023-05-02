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

func JourneysGet(newRequest ReqQueryParameters) (response events.APIGatewayProxyResponse) {
	var err error

	if newRequest.Method != "GET" {
		return createErrorResponse(http.StatusBadRequest, "only GET method is allowed for journeys api")
	}

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

	filter := db.JourneyFilter{}
	if newRequest.ID != "" {
		filter.ID, err = strconv.Atoi(newRequest.ID)
		if err != nil {
			log.Println("Invalid query parameter:", err)
			return createErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v invalid ID request for journeys", filter.ID))
		}
	}

	lastJourney, err := DB.GetLastJourneyId()
	if err != nil {
		log.Println("Error while getting last journey ID:", err)
		return createErrorResponse(http.StatusInternalServerError, "error trying to retreive a journey, please try again later")
	}

	if filter.ID > lastJourney.ID {
		log.Printf("Error while getting batch with starting id ID %v: %v", filter.ID, err)
		return createErrorResponse(http.StatusBadRequest, fmt.Sprintf("%v is an invalid ID request for journeys", filter.ID))
	}

	filter.Limit = 3000
	remainingIds := lastJourney.ID - filter.Limit
	if remainingIds < 3000 {
		filter.Limit = remainingIds
	}

	journeys, err := DB.GetJourneys(filter)
	if err != nil {
		log.Println("Error while getting journeys:", err)
		return createErrorResponse(http.StatusInternalServerError, "error while receiving stations, please try again later")
	}

	var responseJSON []byte
	if journeys != nil {
		responseJSON, err = json.Marshal(journeys)
		if err != nil {
			log.Println("error while marshalling journeys:", err)
			return createErrorResponse(http.StatusInternalServerError, "oops! something went wrong while processing your journeys")
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseJSON),
	}
}
