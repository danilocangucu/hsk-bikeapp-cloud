# Cloud Backend Integration for My Helsinki City Bike Single Page App

This repository contains additional code for my [Helsinki City Bike Single Page App](https://github.com/danilocangucu/hsk-bikeapp-solita), providing cloud backend integration using Amazon Web Services Relational Database Service (AWS RDS) and Amazon Web Services Lambda.

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Importing database to AWS RDS](#importing-database-to-aws-rds)
4. [Go Backend](#go-backend)
5. [Lambda function on AWS](#lambda-function-on-aws)
   1) [Creating the Lambda function](#creating-the-lambda-function)
   2) [Deploying the Lambda function](#deploying-the-lambda-function)
   3) [Adding an API Gateway as a trigger](#adding-an-api-gateway-as-a-trigger)
   4) [Using the Lambda function](#using-the-lambda-function)

## Introduction

The code in this repository serves as a guide for integrating cloud services with the Helsinki City Bike app, as part of my Solita's Dev Academy 2023 pre-assignment's application. If you're not familiar with it, I recommend starting with the linked repository above before continuing.

## Getting Started

To begin working on this project, you'll need the following software installed on your computer:
- [Go](https://golang.org/)
  - [Go-MySQL-Driver Package](https://pkg.go.dev/github.com/go-sql-driver/mysql@v1.7.0#section-readme)
  - [AWS Lambda Go Package](https://github.com/aws/aws-lambda-go)
- [SQLite3](https://www.sqlite.org/index.html)
- [MySQL](https://www.mysql.com/)

## Importing database to AWS RDS

Before getting started with AWS RDS, the first step is to convert the previous database from SQLite to MySQL. Please copy the file `database/hsk-city-bike-app.db` from the previous project to the `database` directory or create the previous database with the command:

```
./database/dbcreate.sh
```

When you have the SQLite database ready, convert it to a MySQL database and import to AWS RDS by using the following commands:
```
sqlite3 database/hsk-city-bike-app.db .dump | mysql -h <hostname> -u <username> -p<password> <database name>
```
This command will drop all tables from the database and import it to AWS RDS. Note that the process might take several hours to complete – it took me around a day!

To check the status of the import process, you can use the following commands in another terminal or prompt window:
1. Log in to the AWS RDS database:
```
mysql -h <hostname> -u <username> -p<password> <database name>
```
2. Check the status of the importing process (and all processes in the database):
```
SHOW PROCESSLIST;
```
Note that you can find the hostname, username, password, and database name in my application documents. If you encounter any issues while using these credentials, please do not hesitate to contact me through email or phone.

## Go Backend

The Go backend for this application consists of two main packages: `handlers` and `database`. The `handlers` package contains the `common.go` file, which defines an APIHandler struct and various functions to handle AWS Lambda requests for the "stations" and "journeys" APIs, with a rate limit of one request per second. The `journeys.go` file processes GET requests for the "journeys" API, while the `stations.go` file handles GET and POST requests for the "stations" API. Both involve connecting to the database, validating input, and returning appropriate responses.

In the `database` package, the `db.go` file provides functions to interact with the database and manage station and journey data. It includes the `Db` struct, representing a database connection, and data structures like `Station`, `Journey`, `StationFilter`, and `JourneyFilter`. Functions handle various operations such as opening and closing connections, fetching data, and adding new records. The file uses `sync.WaitGroup` to validate new stations concurrently for efficiency.

To use the application, ensure you insert the correct credentials in the `OpenDatabase` function located in `database/db.go`. The necessary credentials can be found in the application's documents.

## Lambda function on AWS
This section covers the process of creating, deploying, and using a Lambda function on AWS, including configuring the function, and invoking it with an API Gateway trigger.

### Creating the Lambda function
First, you need to export the Go code provided in this repository to create a building version of the app with the right environment and compress it to a zip file. To do this, you can run the following command (on Mac):

```
env GOOS=linux GOARCH=amd64 go build -o main
zip main.zip main
```

For other operating systems, please refer to the [AWS Lambda for Go Package](https://github.com/aws/aws-lambda-go) or Lambda's [official documentation](https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model.html).

### Deploying the Lambda function

After creating the ZIP file, you can proceed to upload it to AWS Lambda and configure the runtime settings to make the Lambda function operational. Here are the steps to follow:

1. Go to AWS Lambda and create a new function.
2. Upload the "main.zip" file.
3. Access the runtime settings and select "Go 1.x" as the runtime.
4. Specify "main" as the handler function in the "Handler" field.
5. Save the changes to the runtime settings.

Under "Function overview", you can find the URL of the function. For example, that's mine:
```
https://33r4rjpg7j4av5diigehcezcy40wzwqm.lambda-url.eu-north-1.on.aws/
```

### Adding an API Gateway as a trigger

When your Lambda function is ready, you can trigger it with an API Gateway by following these steps:

1. In "Function overview", click on "Add trigger";
2. Choose "API Gateway" under "Trigger configuration" and select "Create an API";
3. For the type, I chose "HTTP API" and for security, I chose "Open".

In the "Triggers" section, you'll find detailed information about the API Gateway, including the API endpoint that you'll use to make requests. Here's my API's endpoint:
```
https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp
```

### Using the Lambda function

Retrieve data by calling the two APIs present in the Lambda function: `stations` and `journeys`. To make a request, include the query parameter `api` and an optional `id` parameter. Additionally, the `stations` API can receive POST requests for adding new stations, as demonstrated in example cases 4 and 5 below.

**Stations API:**

1. Retrieve all stations data:
   ```
   https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=stations
   ```
   Result example of the [request](https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=stations):

   ```json
   [
     {
       "FID": 457 ,
       "ID": 405 ,
       "Nimi": "Jollas" ,
       "Namn": "Jollas" ,
       "Name": "Jollas" ,
       "Osoite": "Jollaksentie 33" ,
       "Adress": "Jollasvägen 33" ,
       "Kaupunki": " " ,
       "Stad": " " ,
       "Operaattor": " " ,
       "Kapasiteet": 16 ,
       "Latitude": 25.0617 ,
       "Longitude": 60.1644 ,
       "JourneysFrom": 661 ,
       "JourneysTo": 825
     },
     ...
     {
       "FID": 1 ,
       "ID": 501 ,
       "Nimi": "Hanasaari" ,
       "Namn": "Hanaholmen" ,
       "Name": "Hanasaari" ,
       "Osoite": "Hanasaarenranta 1" ,
       "Adress": "Hanaholmsstranden 1" ,
       "Kaupunki": "Espoo" ,
       "Stad": "Esbo" ,
       "Operaattor": "CityBike Finland" ,
       "Kapasiteet": 10 ,
       "Latitude": 24.8403 ,
       "Longitude": 60.1658 ,
       "JourneysFrom": 2373 ,
       "JourneysTo": 2442
     }
   ]
   ```

2. Retrieve data for a specific station ID (e.g., ID 11):
   ```
   https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=stations&id=11
   ```
   Example result of the [request](https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=stations&id=11):

   ```json
   {
     "FID": 121,
     "ID": 11,
     "Nimi": "Unioninkatu",
     "Namn": "Unionsgatan",
     "Name": "Unioninkatu",
     "Osoite": "Eteläesplanadi 1",
     "Adress": "Södra esplanaden 1",
     "Kaupunki": " ",
     "Stad": " ",
     "Operaattor": " ",
     "Kapasiteet": 22,
     "Latitude": 24.951,
     "Longitude": 60.1675,
     "JourneysFrom": 10579,
     "JourneysTo": 12368
   }
   ```

3. Attempt to retrieve data for a non-existent station ID (e.g., ID 100000) returns an error:
   ```
   https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=stations&id=100000
   ```
   Example result of the [request](https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=stations&id=100000):

   Error message: `error while getting station ID 100000`

   More detailed error information can be found in AWS CloudWatch logs:
   ```
   2023/05/03 09:17:55 Error while getting station ID 100000: sql: no rows in result set
   ```

4. To add a new station, send a POST request to the following URL:
   ```
   https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=stations
   ```
   
   Include a JSON in the request body, for example:
   ```json
   {
     "ID": 1000,
     "Nimi": "New Finnish Name",
     "Namn": "New Swedish Name",
     "Name": "New English Name",
     "Osoite": "New Finnish Address",
     "Adress": "New Swedish Address",
     "Kaupunki": "New Finnish City",
     "Stad": "New Swedish City",
     "Operaattor": "New Operator",
     "Kapasiteet": 10,
     "Latitude": 60.1698,
     "Longitude": 24.9388,
     "JourneysFrom": 0,
     "JourneysTo": 0
   }
   ```
   
   Example result of the request:
   ```
   new station added successfully!
   ```
   
5. If you try to add an existing station with the same JSON again, you will receive the following error message:

   ```json
   [
     "Station with coordinates (60.169800, 24.938801) already exists",
     "station with Finnish name 'New Finnish Name' already exists",
     "Station with Swedish address 'New Swedish Address' already exists",
     "station with Swedish name 'New Swedish Name' already exists",
     "Station with Finnish address 'New Finnish Address' already exists",
     "Station with English name 'New English Name' already exists"
   ]
   ```
   
Please be aware that the validations presented here are solely based on backend checks from the [Helsinki City Bike Single Page App](https://github.com/danilocangucu/hsk-bikeapp-solita). Frontend validations, such as address validation, are not included in this Cloud integration.


**Journeys API:**

The Journeys API is designed to return a limited batch of journey records at a time. By default, the batch size is set to 3000 journeys. This limit can be found in the `JourneysGet` function within the `handlers/journeys.go` file, where `filter.Limit` is set to 3000.

1. Retrieve a batch of journeys without specifying an ID:
   ```
   https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=journeys
   ```
   This returns a batch of journeys between ID 1 and 3000.

   Example result:

   ```json
   [
     {
       "ID": 1,
       "Departure": "2021-05-01 00:00:11",
       "Return": "2021-05-01 00:04:34",
       "DepartureStationId": 138,
       "DepartureStationName": "Arabiankatu",
       "ReturnStationId": 138,
       "ReturnStationName": "Arabiankatu",
       "CoveredDistanceM": 1057,
       "DurationSec": 259
     },
     ...
     {
       "ID": 3000,
       "Departure": "2021-05-01 11:31:13",
       "Return": "2021-05-01 11:44:24",
       "DepartureStationId": 52,
       "DepartureStationName": "Heikkilänaukio",
       "ReturnStationId": 29,
       "ReturnStationName": "Baana",
       "CoveredDistanceM": 2791,
       "DurationSec": 787
     }
   ]
   ```

2. Retrieve a batch of journeys specifying an ID, for example, 3000:
   ```
   https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=journeys&id=3000
   ```
   This returns a batch of journeys between ID 3001 and 6000.

   Example result:

   ```json
   [
     {
       "ID": 3001,
       "Departure": "2021-05-01 11:31:16",
       "Return": "2021-05-01 11:52:56",
       "DepartureStationId": 529,
       "DepartureStationName": "Keilaniemi (M)",
       "ReturnStationId": 525,
       "ReturnStationName": "Mäntyviita",
       "CoveredDistanceM": 3834,
       "DurationSec": 1295
     },
     ...
     {
       "ID": 6000,
       "Departure": "2021-05-01 15:24:42",
       "Return": "2021-05-01 16:24:49",
       "DepartureStationId": 65,
       "DepartureStationName": "Hernesaarenranta",
       "ReturnStationId": 66,
       "ReturnStationName": "Ehrenströmintie",
       "CoveredDistanceM": 5035,
       "DurationSec": 3603
     }
   ]
   ```

3. Invalid ID example, for example, -3000:
   ```
   https://2b9nuc6zm3.execute-api.eu-north-1.amazonaws.com/default/hsk-bikeapp?api=journeys&id=-3000
   ```
   Error message: `-3000 is an invalid ID`
