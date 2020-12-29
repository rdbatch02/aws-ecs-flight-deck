package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/rdbatch02/ecs-flight-deck/domain"
)

var dynamoService *dynamodb.DynamoDB
var ecsService *ecs.ECS

func main() {
	lambda.Start(HandleRequest)
}

func init() {
	awsSession := session.Must(session.NewSession())
	dynamoService = dynamodb.New(awsSession)
	ecsService = ecs.New(awsSession)
}

func HandleRequest(ctx context.Context) error {
	tableName := os.Getenv("FLIGHT_DECK_TABLE_NAME")
	if len(tableName) == 0 {
		fmt.Printf("FLIGHT_DECK_TABLE_NAME must not be empty")
		return errors.New("FLIGHT_DECK_TABLE_NAME must not be empty")
	}

	clusters, err := ecsService.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		fmt.Printf("Failed to enumerate clusters %s", err)
		return err
	}
	clusterArns := clusters.ClusterArns
	for _, clusterArn := range clusterArns {
		fmt.Printf("Finding services for cluster: %s", *clusterArn)
		services, err := getClusterServices(ecsService, *clusterArn)
		if err != nil {
			return err
		}
		saveServices(services, dynamoService, tableName)
	}
	return nil // No errors!
}

func getClusterServices(ecsService *ecs.ECS, clusterArn string) ([]domain.EcsService, error) {
	listRequest := ecs.ListServicesInput{
		Cluster: &clusterArn,
	}
	serviceList, err := ecsService.ListServices(&listRequest)
	if err != nil {
		fmt.Printf("Failed to list services %s", err)
		return nil, err
	}
	request := ecs.DescribeServicesInput{
		Cluster:  &clusterArn,
		Services: serviceList.ServiceArns,
	}
	// Get services from AWS
	services, err := ecsService.DescribeServices(&request)
	if err != nil {
		fmt.Printf("Failed to enumerate services %s", err)
		return nil, err
	}
	// Convert services to flight-deck model
	fdServices := make([]domain.EcsService, len(services.Services))
	for i, service := range services.Services {
		fdServices[i] = ecsToFdService(service)
	}
	return fdServices, nil
}

func ecsToFdService(service *ecs.Service) domain.EcsService {
	serv := new(domain.EcsService)
	serv.ClusterArn = *service.ClusterArn
	serv.ServiceArn = *service.ServiceArn
	serv.CreatedAt = service.CreatedAt
	serv.LaunchType = *service.LaunchType
	serv.Status = *service.Status
	serv.Name = *service.ServiceName
	serv.DesiredCount = *service.DesiredCount
	serv.RunningCount = *service.RunningCount
	serv.PendingCount = *service.PendingCount
	return *serv
}

func saveServices(services []domain.EcsService, dynamoService *dynamodb.DynamoDB, tableName string) error {
	for _, service := range services {
		ddbService, err := dynamodbattribute.MarshalMap(service)
		if err == nil {
			fmt.Printf("Saving service: %s", service.Name)
			ddbInput := &dynamodb.PutItemInput{
				Item:      ddbService,
				TableName: aws.String(tableName),
			}
			_, err := dynamoService.PutItem(ddbInput)
			if err != nil {
				fmt.Printf("Failed to save service %s", err)
				return err
			}
		}
	}
	return nil
}
