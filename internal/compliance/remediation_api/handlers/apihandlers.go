package apihandlers

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
	"github.com/panther-labs/panther/internal/compliance/remediation_api/remediation"
)

var (
	sqsQueueURL = os.Getenv("SQS_QUEUE_URL")

	awsSession                        = session.Must(session.NewSession())
	sqsClient  sqsiface.SQSAPI        = sqs.New(awsSession)
	invoker    remediation.InvokerAPI = remediation.NewInvoker(session.Must(session.NewSession()))

	//RemediationLambdaNotFound is the Error when the remediation Lambda is not found
	RemediationLambdaNotFound = &models.Error{Message: aws.String("Remediation Lambda not found or misconfigured")}
)

func badRequest(errorMessage *string) *events.APIGatewayProxyResponse {
	errModel := &models.Error{Message: errorMessage}
	body, err := jsoniter.MarshalToString(errModel)
	if err != nil {
		zap.L().Error("errModel.MarshalBinary failed", zap.Error(err))
		body = "invalid request"
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: body}
	}
	return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: body}
}
