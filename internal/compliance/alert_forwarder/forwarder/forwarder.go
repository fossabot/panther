package forwarder

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var (
	alertQueueURL                 = os.Getenv("ALERTING_QUEUE_URL")
	awsSession                    = session.Must(session.NewSession())
	sqsClient     sqsiface.SQSAPI = sqs.New(awsSession)
)

// Handle forwards an alert to the alert delivery SQS queue
func Handle(event *models.Alert) error {
	zap.L().Info("received alert", zap.String("policyId", *event.PolicyID))

	msgBody, err := jsoniter.Marshal(event)
	if err != nil {
		return err
	}
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(alertQueueURL),
		MessageBody: aws.String(string(msgBody)),
	}
	_, err = sqsClient.SendMessage(input)
	if err != nil {
		zap.L().Warn("failed to send message to remediation", zap.Error(err))
		return err
	}
	zap.L().Info("successfully triggered alert action")

	return nil
}
