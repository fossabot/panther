package apihandlers

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
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

type mockSqsClient struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (m *mockSqsClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

var input = &models.RemediateResource{
	PolicyID:   "policyId",
	ResourceID: "resourceId",
}

func TestRemediateResource(t *testing.T) {
	mockInvoker := &mockInvoker{}
	invoker = mockInvoker
	mockSqsClient := &mockSqsClient{}
	sqsClient = mockSqsClient

	serializedPayload, _ := jsoniter.MarshalToString(input)
	request := &events.APIGatewayProxyRequest{Body: serializedPayload}

	mockInvoker.On("Remediate", input).Return(nil)

	response := RemediateResource(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "", response.Body)
	mockInvoker.AssertExpectations(t)
	mockSqsClient.AssertExpectations(t)
}

func TestRemediateResourceMissingParameters(t *testing.T) {
	mockInvoker := &mockInvoker{}
	invoker = mockInvoker
	mockSqsClient := &mockSqsClient{}
	sqsClient = mockSqsClient

	payload := &models.RemediateResource{}
	serializedPayload, _ := payload.MarshalBinary()
	request := &events.APIGatewayProxyRequest{Body: string(serializedPayload)}

	response := RemediateResource(request)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	mockInvoker.AssertExpectations(t)
	mockSqsClient.AssertExpectations(t)
}

func TestRemediateResourceInvalidInput(t *testing.T) {
	mockInvoker := &mockInvoker{}
	invoker = mockInvoker
	mockSqsClient := &mockSqsClient{}
	sqsClient = mockSqsClient

	request := &events.APIGatewayProxyRequest{Body: "not-a-json"}

	response := RemediateResource(request)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	mockInvoker.AssertExpectations(t)
	mockSqsClient.AssertExpectations(t)
}

func TestRemediateResourceAsync(t *testing.T) {
	mockInvoker := &mockInvoker{}
	invoker = mockInvoker
	mockSqsClient := &mockSqsClient{}
	sqsClient = mockSqsClient
	sqsQueueURL = "sqsQueueURL"

	serializedPayload, _ := input.MarshalBinary()
	request := &events.APIGatewayProxyRequest{Body: string(serializedPayload)}

	expectedSqsInput := &sqs.SendMessageInput{
		MessageBody: aws.String(string(serializedPayload)),
		QueueUrl:    aws.String(sqsQueueURL),
	}
	mockSqsClient.On("SendMessage", expectedSqsInput).Return(&sqs.SendMessageOutput{}, nil)

	response := RemediateResourceAsync(request)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "", response.Body)
	mockInvoker.AssertExpectations(t)
	mockSqsClient.AssertExpectations(t)
}

func TestRemediateResourceLambdaDoesntExist(t *testing.T) {
	mockInvoker := &mockInvoker{}
	invoker = mockInvoker

	serializedPayload, _ := jsoniter.MarshalToString(input)
	request := &events.APIGatewayProxyRequest{Body: serializedPayload}

	mockInvoker.On("Remediate", input).Return(
		&genericapi.DoesNotExistError{Message: "there is no aws remediation lambda configured for organization"})
	expectedResponseBody := &models.Error{Message: aws.String("Remediation Lambda not found or misconfigured")}

	response := RemediateResource(request)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	responseBody := &models.Error{}
	assert.NoError(t, jsoniter.UnmarshalFromString(response.Body, responseBody))
	assert.Equal(t, expectedResponseBody, responseBody)

	mockInvoker.AssertExpectations(t)
}
