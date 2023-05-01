package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	LastId int
	Limit  int
}

func OpenDatabase() (Db, error) {

	dbHost := "INSERT-HOST-ADDRESS"
	dbPort := "INSERT-PORT"
	dbUser := "INSERT-USER"
	dbPassword := "INSERT-PASSWOD"
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
	fmt.Println("db opened, returning")

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
		log.Println("db.go error Error while getting station:", err)
		return Station{}, err
	}

	return station, nil
}
