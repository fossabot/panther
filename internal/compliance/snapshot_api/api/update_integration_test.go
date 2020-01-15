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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb/modelstest"
)

func TestUpdateIntegrationSettings(t *testing.T) {
	mockClient := &modelstest.MockDDBClient{}
	db = &ddb.DDB{Client: mockClient, TableName: "test"}

	resp := &dynamodb.UpdateItemOutput{}
	mockClient.On("UpdateItem", mock.Anything).Return(resp, nil)

	result := apiTest.UpdateIntegrationSettings(&models.UpdateIntegrationSettingsInput{
		ScanEnabled:      aws.Bool(false),
		IntegrationID:    aws.String(testIntegrationID),
		AWSAccountID:     aws.String(testAccountID),
		IntegrationLabel: aws.String("NewAWSTestingAccount"),
		ScanIntervalMins: aws.Int(1440),
	})

	assert.NoError(t, result)
	mockClient.AssertExpectations(t)
}

func TestUpdateIntegrationSettingsAwsS3Type(t *testing.T) {
	mockClient := &modelstest.MockDDBClient{}
	db = &ddb.DDB{Client: mockClient, TableName: "test"}

	resp := &dynamodb.UpdateItemOutput{}
	mockClient.On("UpdateItem", mock.Anything).Return(resp, nil)

	result := apiTest.UpdateIntegrationSettings(&models.UpdateIntegrationSettingsInput{
		S3Buckets: aws.StringSlice([]string{"test-bucket-1", "test-bucket-2/*"}),
		KmsKeys:   aws.StringSlice([]string{"arn:aws:kms:us-west-2:415773754570:key/27803c7e-9fa5-4fcb-9525-ee11c953d329"}),
	})

	assert.NoError(t, result)
	mockClient.AssertExpectations(t)
}

func TestUpdateIntegrationValidTime(t *testing.T) {
	now := time.Now()
	validator, err := models.Validator()
	require.NoError(t, err)
	assert.NoError(t, validator.Struct(&models.UpdateIntegrationLastScanStartInput{
		IntegrationID:     aws.String(testIntegrationID),
		ScanStatus:        aws.String(models.StatusOK),
		LastScanStartTime: &now,
	}))
}

func TestUpdateIntegrationLastScanStart(t *testing.T) {
	mockClient := &modelstest.MockDDBClient{}
	db = &ddb.DDB{Client: mockClient, TableName: "test"}

	resp := &dynamodb.UpdateItemOutput{}
	mockClient.On("UpdateItem", mock.Anything).Return(resp, nil)

	lastScanEndTime, err := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	require.NoError(t, err)

	result := apiTest.UpdateIntegrationLastScanStart(&models.UpdateIntegrationLastScanStartInput{
		IntegrationID:     aws.String(testIntegrationID),
		LastScanStartTime: &lastScanEndTime,
		ScanStatus:        aws.String(models.StatusOK),
	})

	assert.NoError(t, result)
	mockClient.AssertExpectations(t)
}

func TestUpdateIntegrationLastScanEnd(t *testing.T) {
	mockClient := &modelstest.MockDDBClient{}
	db = &ddb.DDB{Client: mockClient, TableName: "test"}

	resp := &dynamodb.UpdateItemOutput{}

	update := expression.Set(
		expression.Name("lastScanEndTime"),
		expression.Value("2009-11-10T23:00:00Z"),
	)
	update = update.Set(
		expression.Name("lastScanErrorMessage"),
		expression.Value("something went wrong"),
	)
	update = update.Set(
		expression.Name("scanStatus"),
		expression.Value(models.StatusError),
	)
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	require.NoError(t, err)

	expected := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			"integrationId": {S: aws.String(testIntegrationID)},
		},
		TableName:        aws.String("test"),
		UpdateExpression: expr.Update(),
	}
	mockClient.On("UpdateItem", expected).Return(resp, nil)

	lastScanEndTime, err := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	require.NoError(t, err)

	result := apiTest.UpdateIntegrationLastScanEnd(&models.UpdateIntegrationLastScanEndInput{
		IntegrationID:        aws.String(testIntegrationID),
		LastScanEndTime:      &lastScanEndTime,
		LastScanErrorMessage: aws.String("something went wrong"),
		ScanStatus:           aws.String(models.StatusError),
	})

	assert.NoError(t, result)
	mockClient.AssertExpectations(t)
}
