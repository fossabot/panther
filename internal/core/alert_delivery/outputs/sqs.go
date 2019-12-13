package outputs

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

// Sqs sends an alert to an SQS Queue.
// nolint: dupl
func (client *OutputClient) Sqs(alert *alertmodels.Alert, config *outputmodels.SqsConfig) *AlertDeliveryError {
	outputMessage := &sqsOutputMessage{
		ID:          alert.PolicyID,
		Name:        alert.PolicyName,
		VersionID:   alert.PolicyVersionID,
		Description: alert.PolicyDescription,
		Runbook:     alert.Runbook,
		Severity:    alert.Severity,
		Tags:        alert.Tags,
	}

	serializedMessage, err := jsoniter.MarshalToString(outputMessage)
	if err != nil {
		zap.L().Error("Failed to serialize message", zap.Error(err))
		return &AlertDeliveryError{Message: "Failed to serialize message"}
	}

	sqsSendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    config.QueueURL,
		MessageBody: aws.String(serializedMessage),
	}

	sqsClient, err := client.getSqsClient(*config.QueueURL)
	if err != nil {
		return &AlertDeliveryError{Message: "Failed to create Sqs client for queue", Permanent: true}
	}

	_, err = sqsClient.SendMessage(sqsSendMessageInput)
	if err != nil {
		zap.L().Error("Failed to send message to SQS queue", zap.Error(err))
		return &AlertDeliveryError{Message: "Failed to send message to SQS queue"}
	}
	return nil
}

//sqsOutputMessage contains the fields that will be included in the SQS message
type sqsOutputMessage struct {
	ID          *string   `json:"id"`
	Name        *string   `json:"name,omitempty"`
	VersionID   *string   `json:"versionId,omitempty"`
	Description *string   `json:"description,omitempty"`
	Runbook     *string   `json:"runbook,omitempty"`
	Severity    *string   `json:"severity"`
	Tags        []*string `json:"tags,omitempty"`
}

func (client *OutputClient) getSqsClient(queueURL string) (sqsiface.SQSAPI, error) {
	// Queue URL is like "https://sqs.us-west-2.amazonaws.com/415773754570/panther-alert-queue"
	region := strings.Split(queueURL, ".")[1]
	sqsClient, ok := client.sqsClients[region]
	if !ok {
		sqsClient = sqs.New(client.session, aws.NewConfig().WithRegion(region))
		client.sqsClients[region] = sqsClient
	}
	return sqsClient, nil
}
