package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var ecsService *ecs.ECS
var lambdaService *lambdaSdk.Lambda
var refreshClockLambdaName string
var isLocal bool
var ginLambda *ginadapter.GinLambda

func inLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func init() {
	awsSession := session.Must(session.NewSession())
	ecsService = ecs.New(awsSession)

	refreshClockLambdaName = os.Getenv("REFRESH_CLOCK_LAMBDA_NAME")
	if len(refreshClockLambdaName) == 0 {
		panic(errors.New("REFRESH_CLOCK_LAMBDA_NAME must not be empty"))
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.PUT("/api/service/restart", restartService)
	return router
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if inLambda() {
		fmt.Println("running aws lambda in aws")
		ginLambda = ginadapter.New(setupRouter())
		lambda.Start(Handler)
	} else {
		fmt.Println("running aws lambda in local")
		log.Fatal(http.ListenAndServe(":8080", setupRouter()))
	}
}

type RestartRequest struct {
	ClusterArn string
	ServiceArn string
}

func restartService(ctx *gin.Context) {
	var request RestartRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	fmt.Printf("Restarting service %s", request.ServiceArn)
	_, err := ecsService.UpdateService(&ecs.UpdateServiceInput{
		ForceNewDeployment: aws.Bool(true),
		Service:            &request.ClusterArn,
		Cluster:            &request.ServiceArn,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		go triggerRefresh()
		ctx.Status(http.StatusOK)
	}
}

func triggerRefresh() {
	// Trigger refresh
	lambdaService.Invoke(&lambdaSdk.InvokeInput{
		FunctionName:   &refreshClockLambdaName,
		InvocationType: aws.String("Event"),
	})
}
