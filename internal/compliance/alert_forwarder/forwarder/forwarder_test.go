package forwarder

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

type mockSqsClient struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (m *mockSqsClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

func init() {
	alertQueueURL = "alertQueueURL"
}

func TestHandleAlert(t *testing.T) {
	mockSqsClient := &mockSqsClient{}
	sqsClient = mockSqsClient

	input := &models.Alert{
		PolicyID: aws.String("policyId"),
	}

	expectedMsgBody, _ := json.Marshal(input)
	expectedInput := &sqs.SendMessageInput{
		QueueUrl:    aws.String("alertQueueURL"),
		MessageBody: aws.String(string(expectedMsgBody)),
	}

	mockSqsClient.On("SendMessage", expectedInput).Return(&sqs.SendMessageOutput{}, nil)
	require.NoError(t, Handle(input))
	mockSqsClient.AssertExpectations(t)
}

func TestHandleAlertSqsError(t *testing.T) {
	mockSqsClient := &mockSqsClient{}
	sqsClient = mockSqsClient

	input := &models.Alert{
		PolicyID: aws.String("policyId"),
	}

	mockSqsClient.On("SendMessage", mock.Anything).Return(&sqs.SendMessageOutput{}, errors.New("error"))
	require.Error(t, Handle(input))
	mockSqsClient.AssertExpectations(t)
}
