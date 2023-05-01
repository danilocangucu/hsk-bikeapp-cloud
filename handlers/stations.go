package handlers

import (
	"encoding/json"
	"fmt"
	db "hsk-bikeapp-solita/database"
	"log"
	"strconv"
)

var DB db.Db

func StationsGet(newRequest ReqQueryParameters) (result string, err error) {
	if newRequest.ID < "0" {
		log.Println("invalid query parameter:", newRequest.ID)
		return "", fmt.Errorf("%v is an invalid ID", newRequest.ID)
	}

	// Open database connection
	DB, err = db.OpenDatabase()
	if err != nil {
		log.Println("Error opening database:", err)
		return "", fmt.Errorf("an error has been produced trying to access the database")
	}
	defer DB.CloseDatabase()

	// Create filter based on query parameters
	filter := db.StationFilter{}
	if newRequest.ID != "" {
		filter.ID, err = strconv.Atoi(newRequest.ID)
		if err != nil {
			log.Println("Invalid query parameter:", err)
			return "", fmt.Errorf("%v is an invalid ID", err)
		}
	}

	// Get stations based on filter
	var station db.Station
	var stations []db.Station
	if filter.ID != 0 {
		station, err = DB.GetSingleStation(filter)
		if err != nil {
			log.Printf("Error while getting station ID=%v: %v", filter.ID, err)
			return "", fmt.Errorf("error while getting station ID %v", filter.ID)
		}
	} else if filter.ID == 0 {
		stations, err = DB.GetAllStations()
		if err != nil {
			log.Println("Error while getting stations:", err)
			return "", fmt.Errorf("error while getting stations")
		}
	}

	// Marshal response into JSON
	var responseJSON []byte
	if station != (db.Station{}) {
		responseJSON, err = json.Marshal(station)
		if err != nil {
			log.Println("error while marshalling station:", err)
			return "", fmt.Errorf("oops! something went wrong while processing your station request")
		}
	} else {
		responseJSON, err = json.Marshal(stations)
		if err != nil {
			log.Println("error while marshalling stations:", err)
			return "", fmt.Errorf("oops! something went wrong while processing your stations request")
		}
	}

	result = string(responseJSON)
	return result, nil
}
