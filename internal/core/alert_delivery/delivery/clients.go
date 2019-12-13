package delivery

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/panther-labs/panther/internal/core/alert_delivery/outputs"
)

var (
	awsSession = session.Must(session.NewSession())

	// We will always need the Lambda client (to get output details)
	lambdaClient lambdaiface.LambdaAPI = lambda.New(awsSession)

	outputClient outputs.API = outputs.New(awsSession)

	// Lazy-load the SQS client - we only need it to retry failed alerts
	sqsClient sqsiface.SQSAPI
)

func getSQSClient() sqsiface.SQSAPI {
	if sqsClient == nil {
		sqsClient = sqs.New(awsSession)
	}
	return sqsClient
}
