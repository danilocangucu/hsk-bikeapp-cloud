# Cloud Backend Integration for My Helsinki City Bike Single Page App

This repository contains additional code for my [Helsinki City Bike Single Page App](https://github.com/danilocangucu/hsk-bikeapp-solita), providing cloud backend integration using Amazon Web Services Relational Database Service (AWS RDS) and Amazon Web Services Lambda.

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Importing database to AWS RDS](#importing-database-to-aws-rds)
4. [Lambda function on AWS](#lambda-function-on-aws)
   1) [Creating the Lambda function](#creating-the-lambda-function)
   2) [Deploying the Lambda function](#deploying-the-lambda-function)

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

The Go backend consists of two primary components. The first component is a Go package located at `database/db.go` and is called "database." This package provides various functions to interact with the AWS RDS MySQL database, including establishing and terminating database connections, retrieving one or multiple records from the "stations" table, and filtering results based on different criteria. Additionally, it defines structs such as Station and Journey to represent the data. To achieve these functionalities, the package utilizes the `database/sql` and `Go-MySQL-Driver` packages.

The second component is a Go package called "handlers" located in the `handlers/` directory. The `handlers/common.go` file is designed for use with AWS Lambda. The `HandleRequest` function accepts `APIGatewayProxyRequest` inputs and generates appropriate responses based on the request's parameters.

The `handlers/stations.go` file contains a function that retrieves station data from the database using functions from the "database" package and returns it as JSON. This function also handles errors related to invalid parameters or database connectivity issues.

Please ensure that you insert the correct credentials in the `OpenDatabase` function located at `database/db.go`. The necessary credentials are provided in my application's documents.

## Lambda function on AWS
This section covers creating and deploying a Lambda function on AWS.

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
