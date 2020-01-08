package outputs

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

type mockSqsClient struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (m *mockSqsClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

func TestSendSqs(t *testing.T) {
	client := &mockSqsClient{}
	outputClient := &OutputClient{sqsClients: map[string]sqsiface.SQSAPI{"us-west-2": client}}

	sqsOutputConfig := &outputmodels.SqsConfig{
		QueueURL: aws.String("https://sqs.us-west-2.amazonaws.com/123456789012/test-output"),
	}
	alert := &alertmodels.Alert{
		PolicyName:        aws.String("policyName"),
		PolicyID:          aws.String("policyId"),
		PolicyDescription: aws.String("policyDescription"),
		Severity:          aws.String("severity"),
		Runbook:           aws.String("runbook"),
	}

	expectedSqsMessage := &sqsOutputMessage{
		ID:          alert.PolicyID,
		Name:        alert.PolicyName,
		Description: alert.PolicyDescription,
		Severity:    alert.Severity,
		Runbook:     alert.Runbook,
	}
	expectedSerializedSqsMessage, _ := jsoniter.MarshalToString(expectedSqsMessage)
	expectedSqsSendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    sqsOutputConfig.QueueURL,
		MessageBody: aws.String(expectedSerializedSqsMessage),
	}

	client.On("SendMessage", expectedSqsSendMessageInput).Return(&sqs.SendMessageOutput{}, nil)
	result := outputClient.Sqs(alert, sqsOutputConfig)
	assert.Nil(t, result)
	client.AssertExpectations(t)
}
