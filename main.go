package main

import (
	handlers "hsk-bikeapp-solita-cloud/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	apiHandler := &handlers.APIHandler{}
	lambda.Start(apiHandler.HandleRequest)
}
