package refreshlambda

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/rdbatch02/ecs-flight-deck/domain"
)

func main() {
	awsSession := session.Must(session.NewSession())
	dynamoService := dynamodb.New(awsSession)
	ecsService := ecs.New(awsSession)

	clusters, err := ecsService.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		fmt.Printf("Failed to enumerate clusters %s", err)
	}
	clusterArns := clusters.ClusterArns
	for _, clusterArn := range clusterArns {
		services, err := getClusterServices(ecsService, *clusterArn)
		if err != nil {
			saveServices(services, dynamoService)
		}
	}

}

func getClusterServices(ecsService *ecs.ECS, clusterArn string) ([]domain.EcsService, error) {
	request := ecs.DescribeServicesInput{
		Cluster: &clusterArn,
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
	serv.CreatedAt = service.CreatedAt
	serv.LaunchType = *service.LaunchType
	serv.Status = *service.Status
	serv.DesiredCount = *service.DesiredCount
	serv.RunningCount = *service.RunningCount
	serv.PendingCount = *service.PendingCount
	return *serv
}

func saveServices(services []domain.EcsService, dynamoService *dynamodb.DynamoDB) error {
	for _, service := range services {
		ddbService, err := dynamodbattribute.MarshalMap(service)
		if err == nil {
			ddbInput := &dynamodb.PutItemInput{
				Item:      ddbService,
				TableName: aws.String("testtable"), // TODO: Put proper table name here
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
