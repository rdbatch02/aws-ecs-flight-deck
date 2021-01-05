package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func main() {
	lambda.Start(route)
}

func route(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	httpRequest := req.RequestContext.HTTP

	fmt.Printf("Routing %s", httpRequest.Path)

	switch httpRequest.Path {
	case "/api/clusters/":
		// TODO: Figure out variable routing
	default:
		return clientError(http.StatusNotFound)
	}
}

func restartService(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

}

func serverError(err error) (events.APIGatewayV2HTTPResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}
