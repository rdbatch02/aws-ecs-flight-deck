package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/rdbatch02/ecs-flight-deck/domain"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var dynamoService *dynamodb.DynamoDB
var ecsService *ecs.ECS
var tableName string

func init() {
	tableName = os.Getenv("FLIGHT_DECK_TABLE_NAME")
	if len(tableName) == 0 {
		fmt.Printf("FLIGHT_DECK_TABLE_NAME must not be empty")
		panic("FLIGHT_DECK_TABLE_NAME must not be empty")
	}

	awsSession := session.Must(session.NewSession())

	dynamoService = dynamodb.New(awsSession)
	ecsService = ecs.New(awsSession)
	xray.AWS(dynamoService.Client)
	xray.AWS(ecsService.Client)
}

func main() {
	lambda.Start(route)
}

func route(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	httpRequest := req.RequestContext.HTTP

	fmt.Printf("Routing %s", httpRequest.Path)

	switch httpRequest.Path {
	case "/clusters":
		return clustersHandler(req)
	default:
		return clientError(http.StatusNotFound)
	}
}

func clustersHandler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if arn, hasArn := req.QueryStringParameters["arn"]; hasArn {
		return getClusterServices(arn)
	}
	return getAllClusters()
}

func getAllClusters() (events.APIGatewayV2HTTPResponse, error) {
	clusters, err := ecsService.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		fmt.Printf("Failed to enumerate clusters %s", err)
		return serverError(err)
	}
	clusterArns, err := json.Marshal(clusters.ClusterArns)
	if err != nil {
		return serverError(err)
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(clusterArns),
	}, nil
}

func getClusterServices(clusterArn string) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Printf("Searching for services from cluster: %s", clusterArn)
	query := dynamodb.QueryInput{
		TableName:              &tableName,
		ConsistentRead:         aws.Bool(true),
		KeyConditionExpression: aws.String("ClusterArn = :clusterArn"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":clusterArn": {S: &clusterArn},
		},
	}
	services := []domain.EcsService{}
	err := dynamoService.QueryPages(&query, func(page *dynamodb.QueryOutput, lastPage bool) bool {
		for _, item := range page.Items {
			service := domain.EcsService{}
			err := dynamodbattribute.UnmarshalMap(item, &service)
			if err != nil {
				fmt.Printf("Failed to unmarshal Dynamo record %v", err)
			}
			services = append(services, service)
		}
		return true
	})
	if err != nil {
		return serverError(err)
	}

	servicesJSON, _ := json.Marshal(services)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(servicesJSON),
	}, nil
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
