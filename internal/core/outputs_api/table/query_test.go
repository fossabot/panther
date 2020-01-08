package table

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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

var mockItemMap = map[string]*dynamodb.AttributeValue{
	"outputId": {
		S: aws.String("outputId"),
	},
}

var mockScanItemOutput = &dynamodb.ScanOutput{
	Count: aws.Int64(1),
	Items: []map[string]*dynamodb.AttributeValue{mockItemMap},
}

func TestGetOutputByNameOutputNotFound(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("testTable"), DisplayNameIndex: aws.String("displayIndex")}
	expectedKeyCondition := expression.Key("displayName").Equal(expression.Value(aws.String("displayName")))
	expectedQueryExpression, _ := expression.NewBuilder().
		WithKeyCondition(expectedKeyCondition).
		Build()
	expectedQueryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("testTable"),
		IndexName:                 aws.String("displayIndex"),
		ExpressionAttributeNames:  expectedQueryExpression.Names(),
		ExpressionAttributeValues: expectedQueryExpression.Values(),
		KeyConditionExpression:    expectedQueryExpression.KeyCondition(),
	}
	dynamoResponse := &dynamodb.QueryOutput{Items: make([]map[string]*dynamodb.AttributeValue, 0)}
	dynamoDBClient.On("Query", expectedQueryInput).Return(dynamoResponse, nil)

	result, err := table.GetOutputByName(aws.String("displayName"))
	assert.Nil(t, result)
	assert.NoError(t, err)
	dynamoDBClient.AssertExpectations(t)
}

func TestGetOutputByName(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient}

	dynamoResponse := &dynamodb.QueryOutput{Items: make([]map[string]*dynamodb.AttributeValue, 1)}
	dynamoDBClient.On("Query", mock.Anything).Return(dynamoResponse, nil)

	result, err := table.GetOutputByName(aws.String("displayName"))

	assert.NotNil(t, result)
	assert.NoError(t, err)
	dynamoDBClient.AssertExpectations(t)
}

func TestCheckDuplicateNameServiceError(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient}

	dynamoDBClient.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{}, errors.New("failed"))

	result, err := table.GetOutputByName(aws.String("displayName"))

	require.Nil(t, result)
	assert.Error(t, err)
	dynamoDBClient.AssertExpectations(t)
}

func TestGetOutputs(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("testTable")}
	expectedScanExpression, _ := expression.NewBuilder().Build()
	expectedScanInput := &dynamodb.ScanInput{
		TableName:                 aws.String("testTable"),
		ExpressionAttributeNames:  expectedScanExpression.Names(),
		ExpressionAttributeValues: expectedScanExpression.Values(),
	}

	dynamoDBClient.On("Scan", expectedScanInput).Return(mockScanItemOutput, nil)
	expectedResult := &models.AlertOutputItem{
		OutputID: aws.String("outputId"),
	}

	result, err := table.GetOutputs()
	require.NoError(t, err)
	assert.Equal(t, []*models.AlertOutputItem{expectedResult}, result)

	dynamoDBClient.AssertExpectations(t)
}

func TestGetOutputsPagination(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient}

	dynamoResponseInitial := &dynamodb.ScanOutput{Items: []map[string]*dynamodb.AttributeValue{mockItemMap}, LastEvaluatedKey: mockItemMap}
	dynamoResponseFinal := &dynamodb.ScanOutput{Items: []map[string]*dynamodb.AttributeValue{mockItemMap}, LastEvaluatedKey: nil}

	// Returning a response that contains "LastEvaluatedKey" should force the application to re-submit query
	dynamoDBClient.On("Scan", mock.Anything).Return(dynamoResponseInitial, nil).Twice()
	dynamoDBClient.On("Scan", mock.Anything).Return(dynamoResponseFinal, nil)

	expectedResult := &models.AlertOutputItem{
		OutputID: aws.String("outputId"),
	}

	result, err := table.GetOutputs()

	require.NoError(t, err)
	assert.Equal(t, []*models.AlertOutputItem{expectedResult, expectedResult, expectedResult}, result)
	dynamoDBClient.AssertExpectations(t)
}

func TestGetOutput(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{
		client: dynamoDBClient,
		Name:   aws.String("testTable"),
	}

	expectedGetItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("testTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"outputId": {
				S: aws.String("outputId"),
			},
		},
	}
	mockGetItemOutput := &dynamodb.GetItemOutput{Item: mockItemMap}
	dynamoDBClient.On("GetItem", expectedGetItemInput).Return(mockGetItemOutput, nil)

	expectedResult := &models.AlertOutputItem{
		OutputID: aws.String("outputId"),
	}

	result, err := table.GetOutput(aws.String("outputId"))
	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	dynamoDBClient.AssertExpectations(t)
}

func TestGetOutputNoResult(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{
		client: dynamoDBClient,
		Name:   aws.String("testTable"),
	}

	dynamoDBClient.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{}, nil)

	result, err := table.GetOutput(aws.String("outputId"))
	require.Nil(t, result)
	assert.Error(t, err)
	dynamoDBClient.AssertExpectations(t)
}
