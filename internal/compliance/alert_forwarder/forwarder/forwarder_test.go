package forwarder

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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	jsoniter "github.com/json-iterator/go"
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

	expectedMsgBody, err := jsoniter.MarshalToString(input)
	require.NoError(t, err)
	expectedInput := &sqs.SendMessageInput{
		QueueUrl:    aws.String("alertQueueURL"),
		MessageBody: aws.String(expectedMsgBody),
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
