package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
)

var sqsService *sqs.SQS
var lambdaService *lambdaSdk.Lambda

var refreshLambdaName string
var refreshQueueURL string
var delaySeconds int64

func main() {
	lambda.Start(HandleRequest)
}

func init() {
	awsSession := session.Must(session.NewSession())
	sqsService = sqs.New(awsSession)
	lambdaService = lambdaSdk.New(awsSession)

	refreshLambdaName = os.Getenv("REFRESH_LAMBDA_NAME")
	if len(refreshLambdaName) == 0 {
		panic(errors.New("REFRESH_LAMBDA_NAME must not be empty"))
	}

	refreshQueueURL = os.Getenv("REFRESH_QUEUE_URL")
	if len(refreshQueueURL) == 0 {
		fmt.Printf("REFRESH_QUEUE_URL must not be empty")
		panic("REFRESH_QUEUE_URL must not be empty")
	}

	var err error
	delaySeconds, err = strconv.ParseInt(os.Getenv("DELAY_SECONDS"), 10, 64)
	if err != nil {
		panic("DELAY_SECONDS could not be parsed")
	}

}

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Trigger refresh
	lambdaService.Invoke(&lambdaSdk.InvokeInput{
		FunctionName:   &refreshLambdaName,
		InvocationType: aws.String("Event"),
	})
	// Clear refresh queue
	sqsService.PurgeQueue(&sqs.PurgeQueueInput{
		QueueUrl: &refreshQueueURL,
	})

	// Re-arm next refresh
	sqsService.SendMessage(&sqs.SendMessageInput{
		QueueUrl:     &refreshQueueURL,
		DelaySeconds: &delaySeconds,
		MessageBody:  aws.String("RefreshClock"),
	})

	return nil
}
