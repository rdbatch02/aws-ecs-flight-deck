package refreshlambda

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/rdbatch02/ecs-flight-deck/domain"
)

func main() {
	awsSession := session.Must(session.NewSession())
	dynamoService := dynamodb.New(awsSession)
	ecsService := ecs.New(awsSession)

	clusters, err := ecsService.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		fmt.Printf("Failed to enumerate clusters", err)
	}
	clusterArns := clusters.ClusterArns

}

func getClusterServices(ecsService *ecs.ECS, clusterArn string) {
	request := ecs.DescribeServicesInput{
		Cluster: &clusterArn,
	}
	services, err := ecsService.DescribeServices(&request)
	if err != nil {
		fmt.Printf("Failed to enumerate services", err)
	}
	fdServices := make(map[string]domain.EcsService)

	for _, service := range services.Services {
		fdServices[*service.ServiceArn] = ecsToFdService(service)
	}

}

func ecsToFdService(service *ecs.Service) domain.EcsService {
	serv := new(domain.EcsService)
	serv.CreatedAt = service.CreatedAt
	serv.LaunchType = *service.LaunchType
	serv.Status = *service.Status
	serv.DesiredCount = *service.DesiredCount
	serv.RunningCount = *service.RunningCount
	serv.PendingCount = *service.PendingCount
	return *serv
}
