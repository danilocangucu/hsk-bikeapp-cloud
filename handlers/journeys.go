package handlers

import (
	"encoding/json"
	"fmt"
	db "hsk-bikeapp-solita-cloud/database"
	"log"
	"strconv"
)

func JourneysGet(newRequest ReqQueryParameters) (result string, err error) {
	if newRequest.ID < "0" {
		log.Println("invalid query parameter:", newRequest.ID)
		return "", fmt.Errorf("%v is an invalid ID", newRequest.ID)
	}

	DB, err = db.OpenDatabase()
	if err != nil {
		log.Println("Error opening database:", err)
		return "", fmt.Errorf("an error has been produced trying to access the database")
	}
	defer DB.CloseDatabase()

	filter := db.JourneyFilter{}
	if newRequest.ID != "" {
		filter.ID, err = strconv.Atoi(newRequest.ID)
		if err != nil {
			log.Println("Invalid query parameter:", err)
			return "", fmt.Errorf("%v invalid ID request for journeys", filter.ID)
		}
	}

	lastJourney, err := DB.GetLastJourneyId()
	if err != nil {
		log.Println("Error while getting last journey ID:", err)
		return "", fmt.Errorf("error trying to retreive a journey, please try again later")
	}

	if filter.ID > lastJourney.ID {
		log.Printf("Error while getting batch with starting id ID %v: %v", filter.ID, err)
		return "", fmt.Errorf("%v is an invalid ID request for journeys", filter.ID)
	}

	filter.Limit = 3000
	remainingIds := lastJourney.ID - filter.Limit
	if remainingIds < 3000 {
		filter.Limit = remainingIds
	}

	journeys, err := DB.GetJourneys(filter)
	if err != nil {
		log.Println("Error while getting journeys:", err)
		return "", fmt.Errorf("error while receiving stations, please try again later")
	}

	var responseJSON []byte
	if journeys != nil {
		responseJSON, err = json.Marshal(journeys)
		if err != nil {
			log.Println("error while marshalling journeys:", err)
			return "", fmt.Errorf("oops! something went wrong while processing your journeys")
		}
	}

	result = string(responseJSON)
	return result, nil
}
