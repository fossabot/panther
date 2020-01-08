package delivery

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
	"github.com/aws/aws-sdk-go/service/lambda"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
)

func TestGetOutput(t *testing.T) {
	mockClient := &mockLambdaClient{}
	lambdaClient = mockClient
	lambdaResponse := &lambda.InvokeOutput{
		Payload: []byte(`{"displayName": "alert-channel", "outputConfig" : {"slack": {"webhookURL": "slack.com"}}}`),
	}

	mockClient.On("Invoke", mock.Anything).Return(lambdaResponse, nil)
	result, err := getOutput("test-output-id")

	require.Nil(t, err)
	assert.Equal(t, aws.String("alert-channel"), result.DisplayName)
	assert.NotNil(t, result.OutputConfig.Slack)

	// Now the result should be cached
	cachedResult, err := getOutput("test-output-id")

	require.NoError(t, err)
	assert.Equal(t, result, cachedResult)
	mockClient.AssertExpectations(t)
}

func TestGetOutputError(t *testing.T) {
	mockClient := &mockLambdaClient{}
	lambdaClient = mockClient
	mockClient.On("Invoke", mock.Anything).Return((*lambda.InvokeOutput)(nil), errors.New("error"))

	result, err := getOutput("other")
	require.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestGetAlertOutputIds(t *testing.T) {
	mockClient := &mockLambdaClient{}
	lambdaClient = mockClient

	output := &outputmodels.GetDefaultOutputsOutput{
		Defaults: []*outputmodels.DefaultOutputs{
			{
				Severity:  aws.String("INFO"),
				OutputIDs: aws.StringSlice([]string{"default-info-1", "default-info-2"}),
			},
			{
				Severity:  aws.String("MEDIUM"),
				OutputIDs: aws.StringSlice([]string{"default-medium"}),
			},
		},
	}
	payload, err := jsoniter.Marshal(output)
	require.NoError(t, err)
	mockLambdaResponse := &lambda.InvokeOutput{Payload: payload}

	defaultOutputIDsCache = nil // Clear the cache
	mockClient.On("Invoke", mock.Anything).Return(mockLambdaResponse, nil)
	alert := sampleAlert()
	alert.OutputIDs = nil

	result, err := getAlertOutputIds(alert)

	require.NoError(t, err)
	assert.Equal(t, aws.StringSlice([]string{"default-info-1", "default-info-2"}), result)

	// Now the result should be cached
	require.NotNil(t, defaultOutputIDsCache)
	assert.Equal(t, map[string][]*string{
		"INFO":   aws.StringSlice([]string{"default-info-1", "default-info-2"}),
		"MEDIUM": aws.StringSlice([]string{"default-medium"}),
	}, defaultOutputIDsCache.Outputs)

	cachedResult, err := getAlertOutputIds(alert)

	require.NoError(t, err)
	assert.Equal(t, result, cachedResult)
	mockClient.AssertExpectations(t)
}

func TestGetAlertOutputsIdsError(t *testing.T) {
	mockClient := &mockLambdaClient{}
	lambdaClient = mockClient
	mockClient.On("Invoke", mock.Anything).Return((*lambda.InvokeOutput)(nil), errors.New("error"))

	result, err := getOutput("other")
	require.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}
