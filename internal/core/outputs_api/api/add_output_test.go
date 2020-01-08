package api

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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

func TestAddOutputSameNameAlreadyExists(t *testing.T) {
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputTable := &mockOutputTable{}
	outputsTable = mockOutputTable
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification

	mockOutputTable.On("GetOutputByName", aws.String("my-channel")).Return(&models.AlertOutputItem{}, nil)

	input := &models.AddOutputInput{
		DisplayName:  aws.String("my-channel"),
		UserID:       aws.String("userId"),
		OutputConfig: &models.OutputConfig{Slack: &models.SlackConfig{WebhookURL: aws.String("hooks.slack.com")}},
	}

	result, err := (API{}).AddOutput(input)
	require.Nil(t, result)
	assert.Error(t, err)
	mockOutputTable.AssertExpectations(t)
	mockEncryptionKey.AssertExpectations(t)
	mockOutputVerification.AssertExpectations(t)
}

func TestAddOutputPutOutputError(t *testing.T) {
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputTable := &mockOutputTable{}
	outputsTable = mockOutputTable
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification

	mockOutputTable.On("GetOutputByName", aws.String("my-channel")).Return(nil, nil)
	mockOutputTable.On("PutOutput", mock.Anything).Return(errors.New("internal error"))
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockOutputVerification.On("GetVerificationStatus", mock.Anything).Return(aws.String(models.VerificationStatusSuccess), nil)

	input := &models.AddOutputInput{
		UserID:       aws.String("userId"),
		DisplayName:  aws.String("my-channel"),
		OutputConfig: &models.OutputConfig{Slack: &models.SlackConfig{WebhookURL: aws.String("hooks.slack.com")}},
	}

	result, err := (API{}).AddOutput(input)
	assert.Nil(t, result)
	assert.Error(t, err)

	mockOutputTable.AssertExpectations(t)
	mockEncryptionKey.AssertExpectations(t)
	mockOutputVerification.AssertExpectations(t)
}

func TestAddOutputSlack(t *testing.T) {
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputTable := &mockOutputTable{}
	outputsTable = mockOutputTable
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	mockOutputTable.On("GetOutputByName", aws.String("my-channel")).Return(nil, nil)
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockOutputTable.On("PutOutput", mock.Anything).Return(nil)
	mockOutputVerification.On("GetVerificationStatus", mock.Anything).Return(aws.String(models.VerificationStatusSuccess), nil)
	mockDefaultsTable.On("GetDefault", mock.Anything).Return(&models.DefaultOutputsItem{}, nil)
	mockDefaultsTable.On("PutDefaults", mock.Anything).Return(nil)

	input := &models.AddOutputInput{
		UserID:             aws.String("userId"),
		DisplayName:        aws.String("my-channel"),
		OutputConfig:       &models.OutputConfig{Slack: &models.SlackConfig{WebhookURL: aws.String("hooks.slack.com")}},
		DefaultForSeverity: aws.StringSlice([]string{"CRITICAL", "HIGH"}),
	}

	result, err := (API{}).AddOutput(input)
	require.NoError(t, err)

	expected := &models.AddOutputOutput{
		DisplayName:        aws.String("my-channel"),
		OutputType:         aws.String("slack"),
		LastModifiedBy:     aws.String("userId"),
		CreatedBy:          aws.String("userId"),
		OutputConfig:       &models.OutputConfig{Slack: &models.SlackConfig{WebhookURL: aws.String("hooks.slack.com")}},
		OutputID:           result.OutputID,
		CreationTime:       result.CreationTime,
		LastModifiedTime:   result.LastModifiedTime,
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
		DefaultForSeverity: aws.StringSlice([]string{"CRITICAL", "HIGH"}),
	}
	assert.Equal(t, expected, result)

	_, err = uuid.Parse(*result.OutputID)
	assert.NoError(t, err)

	mockOutputTable.AssertExpectations(t)
	mockEncryptionKey.AssertExpectations(t)
	mockOutputVerification.AssertExpectations(t)
}

func TestAddOutputSns(t *testing.T) {
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputTable := &mockOutputTable{}
	outputsTable = mockOutputTable
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification

	mockOutputTable.On("GetOutputByName", aws.String("my-topic")).Return(nil, nil)
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockOutputTable.On("PutOutput", mock.Anything).Return(nil)
	mockOutputVerification.On("GetVerificationStatus", mock.Anything).Return(aws.String(models.VerificationStatusSuccess), nil)

	input := &models.AddOutputInput{
		UserID:       aws.String("userId"),
		DisplayName:  aws.String("my-topic"),
		OutputConfig: &models.OutputConfig{Sns: &models.SnsConfig{TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:MyTopic")}},
	}

	result, err := (API{}).AddOutput(input)
	require.NoError(t, err)

	expected := &models.AddOutputOutput{
		DisplayName:        aws.String("my-topic"),
		OutputType:         aws.String("sns"),
		LastModifiedBy:     aws.String("userId"),
		CreatedBy:          aws.String("userId"),
		OutputConfig:       &models.OutputConfig{Sns: &models.SnsConfig{TopicArn: aws.String("arn:aws:sns:us-west-2:123456789012:MyTopic")}},
		OutputID:           result.OutputID,
		CreationTime:       result.CreationTime,
		LastModifiedTime:   result.LastModifiedTime,
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
	}
	assert.Equal(t, expected, result)

	_, err = uuid.Parse(*result.OutputID)
	assert.NoError(t, err)
}

func TestAddOutputPagerDuty(t *testing.T) {
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputTable := &mockOutputTable{}
	outputsTable = mockOutputTable
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification

	mockOutputTable.On("GetOutputByName", aws.String("my-pagerduty-integration")).Return(nil, nil)
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockOutputTable.On("PutOutput", mock.Anything).Return(nil)
	mockOutputVerification.On("GetVerificationStatus", mock.Anything).Return(aws.String(models.VerificationStatusSuccess), nil)

	input := &models.AddOutputInput{
		UserID:       aws.String("userId"),
		DisplayName:  aws.String("my-pagerduty-integration"),
		OutputConfig: &models.OutputConfig{PagerDuty: &models.PagerDutyConfig{IntegrationKey: aws.String("93ee508cbfea4604afe1c77c2d9b5bbd")}},
	}

	result, err := (API{}).AddOutput(input)
	require.NoError(t, err)

	expected := &models.AddOutputOutput{
		DisplayName:    aws.String("my-pagerduty-integration"),
		OutputType:     aws.String("pagerduty"),
		LastModifiedBy: aws.String("userId"),
		CreatedBy:      aws.String("userId"),
		OutputConfig: &models.OutputConfig{
			PagerDuty: &models.PagerDutyConfig{
				IntegrationKey: aws.String("93ee508cbfea4604afe1c77c2d9b5bbd"),
			},
		},
		OutputID:           result.OutputID,
		CreationTime:       result.CreationTime,
		LastModifiedTime:   result.LastModifiedTime,
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
	}
	assert.Equal(t, expected, result)

	_, err = uuid.Parse(*result.OutputID)
	assert.NoError(t, err)
}

func TestAddOutputEmail(t *testing.T) {
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputTable := &mockOutputTable{}
	outputsTable = mockOutputTable
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification

	mockOutputTable.On("GetOutputByName", aws.String("my-email")).Return(nil, nil)
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockOutputTable.On("PutOutput", mock.Anything).Return(nil)
	mockOutputVerification.On("GetVerificationStatus", mock.Anything).Return(aws.String(models.VerificationStatusSuccess), nil)

	input := &models.AddOutputInput{
		UserID:       aws.String("userId"),
		DisplayName:  aws.String("my-email"),
		OutputConfig: &models.OutputConfig{Email: &models.EmailConfig{DestinationAddress: aws.String("test@test.com")}},
	}

	result, err := (API{}).AddOutput(input)
	require.NoError(t, err)

	expected := &models.AddOutputOutput{
		DisplayName:        aws.String("my-email"),
		OutputType:         aws.String("email"),
		LastModifiedBy:     aws.String("userId"),
		CreatedBy:          aws.String("userId"),
		OutputConfig:       &models.OutputConfig{Email: &models.EmailConfig{DestinationAddress: aws.String("test@test.com")}},
		OutputID:           result.OutputID,
		CreationTime:       result.CreationTime,
		LastModifiedTime:   result.LastModifiedTime,
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
	}
	assert.Equal(t, expected, result)

	_, err = uuid.Parse(*result.OutputID)
	assert.NoError(t, err)
}

func TestAddOutputSqs(t *testing.T) {
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputTable := &mockOutputTable{}
	outputsTable = mockOutputTable
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification

	mockOutputTable.On("GetOutputByName", aws.String("my-queue")).Return(nil, nil)
	mockEncryptionKey.On("EncryptConfig", mock.Anything).Return(make([]byte, 1), nil)
	mockOutputTable.On("PutOutput", mock.Anything).Return(nil)
	mockOutputVerification.On("GetVerificationStatus", mock.Anything).Return(aws.String(models.VerificationStatusSuccess), nil)

	input := &models.AddOutputInput{
		UserID:      aws.String("userId"),
		DisplayName: aws.String("my-queue"),
		OutputConfig: &models.OutputConfig{
			Sqs: &models.SqsConfig{
				QueueURL: aws.String("https://sqs.us-west-2.amazonaws.com/123456789012/test-output"),
			},
		},
	}

	result, err := (API{}).AddOutput(input)
	require.NoError(t, err)

	expected := &models.AddOutputOutput{
		DisplayName:    aws.String("my-queue"),
		OutputType:     aws.String("sqs"),
		LastModifiedBy: aws.String("userId"),
		CreatedBy:      aws.String("userId"),
		OutputConfig: &models.OutputConfig{
			Sqs: &models.SqsConfig{
				QueueURL: aws.String("https://sqs.us-west-2.amazonaws.com/123456789012/test-output"),
			},
		},
		OutputID:           result.OutputID,
		CreationTime:       result.CreationTime,
		LastModifiedTime:   result.LastModifiedTime,
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
	}
	assert.Equal(t, expected, result)

	_, err = uuid.Parse(*result.OutputID)
	assert.NoError(t, err)
}
