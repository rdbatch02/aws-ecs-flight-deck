package main

import (
	"context"
	"errors"
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
var ginLambda *ginadapter.GinLambda

func init() {
	awsSession := session.Must(session.NewSession())
	ecsService = ecs.New(awsSession)

	refreshClockLambdaName = os.Getenv("REFRESH_CLOCK_LAMBDA_NAME")
	if len(refreshClockLambdaName) == 0 {
		panic(errors.New("REFRESH_LAMBDA_NAME must not be empty"))
	}

	router := gin.Default()
	router.PUT("/api/clusters/:clusterArn/services/:serviceArn", restartService)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

func restartService(ctx *gin.Context) {
	clusterArn := ctx.Param("clusterArn")
	serviceArn := ctx.Param("serviceArn")
	_, err := ecsService.UpdateService(&ecs.UpdateServiceInput{
		ForceNewDeployment: aws.Bool(true),
		Service:            &serviceArn,
		Cluster:            &clusterArn,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		triggerRefresh()
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
