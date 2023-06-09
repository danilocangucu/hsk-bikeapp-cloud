package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	connection *sql.DB
}

type Station struct {
	FID          int
	ID           int
	Nimi         string
	Namn         string
	Name         string
	Osoite       string
	Adress       string
	Kaupunki     string
	Stad         string
	Operaattor   string
	Kapasiteet   int
	Latitude     float32
	Longitude    float32
	JourneysFrom int
	JourneysTo   int
}

type Journey struct {
	ID                   int
	Departure            string
	Return               string
	DepartureStationId   int
	DepartureStationName string
	ReturnStationId      int
	ReturnStationName    string
	CoveredDistanceM     float64
	DurationSec          int
}

type StationFilter struct {
	ID        int
	Nimi      string
	Namn      string
	Name      string
	Osoite    string
	Adress    string
	Latitude  float32
	Longitude float32
}

type JourneyFilter struct {
	ID    int
	Limit int
}

func OpenDatabase() (Db, error) {
	dbHost := "INSERT-HOST-ADDRESS"
	dbPort := "INSERT-PORT"
	dbUser := "INSERT-USER"
	dbPassword := "INSERT-PASSWORD"
	dbName := "INSERT-DATABASE-NAME"

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		return Db{}, err
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return Db{connection: db}, nil
}

func (db *Db) CloseDatabase() {
	db.connection.Close()
}

func (db *Db) GetAllStations() (stations []Station, err error) {
	var station Station

	query := "SELECT FID,ID,Nimi,Namn,Name,Osoite,Adress,Kaupunki,Stad,Operaattor,Kapasiteet,x,y,JourneysFrom,JourneysTo from stations ORDER BY FID DESC"

	rows, err := db.connection.Query(query)
	if err != nil {
		return []Station{}, err
	}

	for rows.Next() {
		err := rows.Scan(&station.FID, &station.ID, &station.Nimi, &station.Namn, &station.Name, &station.Osoite, &station.Adress, &station.Kaupunki, &station.Stad, &station.Operaattor, &station.Kapasiteet, &station.Latitude, &station.Longitude, &station.JourneysFrom, &station.JourneysTo)
		if err != nil {
			return []Station{}, err
		}
		stations = append(stations, station)
	}

	defer rows.Close()
	return stations, err
}

func (db *Db) GetSingleStation(filter StationFilter) (station Station, err error) {
	query := "SELECT FID, ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, x, y, JourneysFrom, JourneysTo FROM test.stations WHERE "
	var args []interface{}
	if filter.ID != 0 {
		query += "ID = ?;"
		args = append(args, filter.ID)
	} else {
		conditions := []string{}
		if filter.Nimi != "" {
			conditions = append(conditions, "Nimi = ?")
			args = append(args, filter.Nimi)
		}
		if filter.Namn != "" {
			conditions = append(conditions, "Namn = ?")
			args = append(args, filter.Namn)
		}
		if filter.Name != "" {
			conditions = append(conditions, "Name = ?")
			args = append(args, filter.Name)
		}
		if filter.Osoite != "" {
			conditions = append(conditions, "Osoite = ?")
			args = append(args, filter.Osoite)
		}
		if filter.Adress != "" {
			conditions = append(conditions, "Adress = ?")
			args = append(args, filter.Adress)
		}
		if filter.Latitude != 0 {
			conditions = append(conditions, "x = ?")
			args = append(args, filter.Latitude)
		}
		if filter.Longitude != 0 {
			conditions = append(conditions, "y = ?")
			args = append(args, filter.Longitude)
		}
		query += strings.Join(conditions, " AND ")
	}

	err = db.connection.QueryRow(query, args...).Scan(&station.FID, &station.ID, &station.Nimi, &station.Namn, &station.Name, &station.Osoite, &station.Adress, &station.Kaupunki, &station.Stad, &station.Operaattor, &station.Kapasiteet, &station.Latitude, &station.Longitude, &station.JourneysFrom, &station.JourneysTo)
	if err != nil {
		log.Println("Error while getting station:", err)
		return Station{}, err
	}

	return station, nil
}

func (db *Db) GetLastJourneyId() (lastJourney JourneyFilter, err error) {
	row := db.connection.QueryRow("SELECT MAX(id) FROM all_journeys;")
	err = row.Scan(&lastJourney.ID)
	if err != nil {
		return lastJourney, err
	}
	return lastJourney, nil
}

func (db *Db) GetJourneys(filter JourneyFilter) (journeys []Journey, err error) {
	var journey Journey

	query := fmt.Sprintf("SELECT * FROM all_journeys WHERE id > %v ORDER BY id LIMIT %v", filter.ID, filter.Limit)
	rows, err := db.connection.Query(query)

	if err != nil {
		return journeys, err
	}

	for rows.Next() {
		err := rows.Scan(&journey.ID, &journey.Departure, &journey.Return, &journey.DepartureStationId, &journey.DepartureStationName, &journey.ReturnStationId, &journey.ReturnStationName, &journey.CoveredDistanceM, &journey.DurationSec)
		if err != nil {
			return journeys, err
		}
		journeys = append(journeys, journey)
	}

	defer rows.Close()
	return journeys, err
}

func (db *Db) ValidateNewStation(newStation Station) (errors []error) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{ID: newStation.ID}); err == nil {
			log.Printf("error validating station id %v: %v", newStation.ID, err)
			errors = append(errors, fmt.Errorf("station with ID %d already exists", newStation.ID))
		}
	}()

	wg.Add(3)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Nimi: newStation.Nimi}); err == nil {
			log.Printf("error validating station Finnish name %v: %v", newStation.Nimi, err)
			errors = append(errors, fmt.Errorf("station with Finnish name '%s' already exists", newStation.Nimi))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Namn: newStation.Namn}); err == nil {
			log.Printf("error validating station Swedish name %v: %v", newStation.Namn, err)
			errors = append(errors, fmt.Errorf("station with Swedish name '%s' already exists", newStation.Namn))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Name: newStation.Name}); err == nil {
			log.Printf("error validating station English name %v: %v", newStation.Name, err)
			errors = append(errors, fmt.Errorf("Station with English name '%s' already exists", newStation.Name))
		}
	}()

	wg.Add(2)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Osoite: newStation.Osoite}); err == nil {
			log.Printf("error validating station Finnish address %v: %v", newStation.Osoite, err)
			errors = append(errors, fmt.Errorf("Station with Finnish address '%s' already exists", newStation.Osoite))
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Adress: newStation.Adress}); err == nil {
			log.Printf("error validating station Swedish address %v: %v", newStation.Adress, err)
			errors = append(errors, fmt.Errorf("Station with Swedish address '%s' already exists", newStation.Adress))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := db.GetSingleStation(StationFilter{Latitude: newStation.Latitude, Longitude: newStation.Longitude}); err == nil {
			log.Printf("error validating station coordinates (%f, %f): %v", newStation.Latitude, newStation.Longitude, err)
			errors = append(errors, fmt.Errorf("Station with coordinates (%f, %f) already exists", newStation.Latitude, newStation.Longitude))
		}
	}()

	wg.Wait()

	return errors
}

func (db *Db) AddNewStation(newStation Station) error {
	query := `INSERT INTO stations (ID, Nimi, Namn, Name, Osoite, Adress, Kaupunki, Stad, Operaattor, Kapasiteet, x, y, JourneysFrom, JourneysTo) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.connection.Exec(query, newStation.ID, newStation.Nimi, newStation.Namn, newStation.Name, newStation.Osoite, newStation.Adress, newStation.Kaupunki, newStation.Stad, newStation.Operaattor, newStation.Kapasiteet, newStation.Latitude, newStation.Longitude, newStation.JourneysFrom, newStation.JourneysTo)

	return err
}
