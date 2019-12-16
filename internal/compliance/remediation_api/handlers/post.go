package apihandlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// RemediateResource remediates a resource synchronously
func RemediateResource(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	remediateResource, errorResponse := checkRequest(request)
	if errorResponse != nil {
		return errorResponse
	}

	zap.L().Info("invoking remediation synchronously")

	if err := invoker.Remediate(remediateResource); err != nil {
		if _, ok := err.(*genericapi.DoesNotExistError); ok {
			return gatewayapi.MarshalResponse(RemediationLambdaNotFound, http.StatusNotFound)
		}
		zap.L().Warn("failed to invoke remediation", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	zap.L().Info("successfully invoked remediation",
		zap.Any("policyId", remediateResource.PolicyID),
		zap.Any("resourceId", remediateResource.ResourceID))
	return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

// RemediateResourceAsync triggers remediation for a resource. The remediation is asynchronous
// so the method will return before the resource has been fixed, independently if it was
// successful or failed.
func RemediateResourceAsync(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	remediateResource, errorResponse := checkRequest(request)
	if errorResponse != nil {
		return errorResponse
	}

	zap.L().Info("sending SQS message to trigger asynchronous remediation")

	sendMessageRequest := &sqs.SendMessageInput{
		MessageBody: aws.String(request.Body),
		QueueUrl:    aws.String(sqsQueueURL),
	}

	if _, err := sqsClient.SendMessage(sendMessageRequest); err != nil {
		zap.L().Warn("failed to send message", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	zap.L().Info("successfully triggered asynchronous remediation",
		zap.Any("policyId", remediateResource.PolicyID),
		zap.Any("resourceId", remediateResource.ResourceID))
	return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

func checkRequest(request *events.APIGatewayProxyRequest) (*models.RemediateResource, *events.APIGatewayProxyResponse) {
	var remediateResource models.RemediateResource

	if err := jsoniter.UnmarshalFromString(request.Body, &remediateResource); err != nil {
		return nil, badRequest(aws.String("invalid request"))
	}

	if err := remediateResource.Validate(nil); err != nil {
		return nil, badRequest(aws.String(err.Error()))
	}
	return &remediateResource, nil
}
