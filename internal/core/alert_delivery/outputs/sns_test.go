package outputs

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

type mockSnsClient struct {
	snsiface.SNSAPI
	mock.Mock
}

func (m *mockSnsClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestSendSns(t *testing.T) {
	client := &mockSnsClient{}
	outputClient := &OutputClient{snsClients: map[string]snsiface.SNSAPI{"us-west-2": client}}

	snsOutputConfig := &outputmodels.SnsConfig{
		TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:test-sns-output"),
	}
	alert := &alertmodels.Alert{
		PolicyName:        aws.String("policyName"),
		PolicyID:          aws.String("policyId"),
		PolicyDescription: aws.String("policyDescription"),
		Severity:          aws.String("severity"),
		Runbook:           aws.String("runbook"),
	}

	expectedSnsMessage := &snsOutputMessage{
		ID:          alert.PolicyID,
		Name:        alert.PolicyName,
		Description: alert.PolicyDescription,
		Severity:    alert.Severity,
		Runbook:     alert.Runbook,
	}
	expectedSerializedSnsMessage, _ := jsoniter.MarshalToString(expectedSnsMessage)
	expectedSnsPublishInput := &sns.PublishInput{
		TopicArn: snsOutputConfig.TopicArn,
		Message:  aws.String(expectedSerializedSnsMessage),
	}

	client.On("Publish", expectedSnsPublishInput).Return(&sns.PublishOutput{}, nil)
	result := outputClient.Sns(alert, snsOutputConfig)
	assert.Nil(t, result)
	client.AssertExpectations(t)
}
