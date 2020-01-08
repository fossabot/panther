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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var mockUpdateItemOutput = &dynamodb.UpdateItemOutput{
	Attributes: map[string]*dynamodb.AttributeValue{
		"outputId": {
			S: aws.String("outputId"),
		},
	},
}
var mockUpdateItemAlertOutput = &models.AlertOutputItem{
	OutputID:           aws.String("outputId"),
	DisplayName:        aws.String("displayName"),
	LastModifiedBy:     aws.String("lastModifiedBy"),
	LastModifiedTime:   aws.String("lastModifiedTime"),
	OutputType:         aws.String("outputType"),
	VerificationStatus: aws.String("verificationStatus"),
	EncryptedConfig:    make([]byte, 1),
}

func TestUpdateOutput(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}

	expectedUpdateExpression := expression.
		Set(expression.Name("displayName"), expression.Value(mockUpdateItemAlertOutput.DisplayName)).
		Set(expression.Name("lastModifiedBy"), expression.Value(mockUpdateItemAlertOutput.LastModifiedBy)).
		Set(expression.Name("lastModifiedTime"), expression.Value(mockUpdateItemAlertOutput.LastModifiedTime)).
		Set(expression.Name("outputType"), expression.Value(mockUpdateItemAlertOutput.OutputType)).
		Set(expression.Name("encryptedConfig"), expression.Value(mockUpdateItemAlertOutput.EncryptedConfig)).
		Set(expression.Name("verificationStatus"), expression.Value(mockUpdateItemAlertOutput.VerificationStatus))

	expectedConditionExpression := expression.Name("outputId").Equal(expression.Value(mockUpdateItemAlertOutput.OutputID))

	expectedExpression, _ := expression.NewBuilder().
		WithCondition(expectedConditionExpression).
		WithUpdate(expectedUpdateExpression).
		Build()

	expectedUpdateItemInput := &dynamodb.UpdateItemInput{
		Key: DynamoItem{
			"outputId": {S: aws.String("outputId")},
		},
		TableName:                 aws.String("TableName"),
		UpdateExpression:          expectedExpression.Update(),
		ConditionExpression:       expectedExpression.Condition(),
		ExpressionAttributeNames:  expectedExpression.Names(),
		ExpressionAttributeValues: expectedExpression.Values(),
		ReturnValues:              aws.String(dynamodb.ReturnValueAllNew),
	}
	expectedResult := &models.AlertOutputItem{
		OutputID: aws.String("outputId"),
	}

	dynamoDBClient.On("UpdateItem", expectedUpdateItemInput).Return(mockUpdateItemOutput, nil)
	result, err := table.UpdateOutput(mockUpdateItemAlertOutput)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	dynamoDBClient.AssertExpectations(t)
}

func TestUpdateOutputWithoutVerificationStatus(t *testing.T) {
	var mockUpdateItemAlertOutput = &models.AlertOutputItem{
		OutputID:         aws.String("outputId"),
		DisplayName:      aws.String("displayName"),
		LastModifiedBy:   aws.String("lastModifiedBy"),
		LastModifiedTime: aws.String("lastModifiedTime"),
		OutputType:       aws.String("outputType"),
		EncryptedConfig:  make([]byte, 1),
	}

	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}

	expectedUpdateExpression := expression.
		Set(expression.Name("displayName"), expression.Value(mockUpdateItemAlertOutput.DisplayName)).
		Set(expression.Name("lastModifiedBy"), expression.Value(mockUpdateItemAlertOutput.LastModifiedBy)).
		Set(expression.Name("lastModifiedTime"), expression.Value(mockUpdateItemAlertOutput.LastModifiedTime)).
		Set(expression.Name("outputType"), expression.Value(mockUpdateItemAlertOutput.OutputType)).
		Set(expression.Name("encryptedConfig"), expression.Value(mockUpdateItemAlertOutput.EncryptedConfig))

	expectedConditionExpression := expression.Name("outputId").Equal(expression.Value(mockUpdateItemAlertOutput.OutputID))

	expectedExpression, _ := expression.NewBuilder().
		WithCondition(expectedConditionExpression).
		WithUpdate(expectedUpdateExpression).
		Build()

	expectedUpdateItemInput := &dynamodb.UpdateItemInput{
		Key: DynamoItem{
			"outputId": {S: aws.String("outputId")},
		},
		TableName:                 aws.String("TableName"),
		UpdateExpression:          expectedExpression.Update(),
		ConditionExpression:       expectedExpression.Condition(),
		ExpressionAttributeNames:  expectedExpression.Names(),
		ExpressionAttributeValues: expectedExpression.Values(),
		ReturnValues:              aws.String(dynamodb.ReturnValueAllNew),
	}
	expectedResult := &models.AlertOutputItem{
		OutputID: aws.String("outputId"),
	}

	dynamoDBClient.On("UpdateItem", expectedUpdateItemInput).Return(mockUpdateItemOutput, nil)
	result, err := table.UpdateOutput(mockUpdateItemAlertOutput)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	dynamoDBClient.AssertExpectations(t)
}

func TestUpdateOutputDoesNotExist(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}

	dynamoDBClient.On("UpdateItem", mock.Anything).Return(
		mockUpdateItemOutput,
		awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "attribute does not exist", nil))

	result, error := table.UpdateOutput(mockUpdateItemAlertOutput)
	assert.Nil(t, result)
	assert.Error(t, error)
	assert.NotNil(t, error.(*genericapi.DoesNotExistError))
	dynamoDBClient.AssertExpectations(t)
}

func TestUpdateOutputServiceError(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}

	dynamoDBClient.On("UpdateItem", mock.Anything).Return(
		mockUpdateItemOutput,
		awserr.New(dynamodb.ErrCodeResourceNotFoundException, "table does not exist", nil))

	result, err := table.UpdateOutput(mockUpdateItemAlertOutput)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NotNil(t, err.(*genericapi.AWSError))
	dynamoDBClient.AssertExpectations(t)
}

func TestUpdateMarshallingError(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}
	mockUpdateItemOutput.Attributes["outputId"] = &dynamodb.AttributeValue{BOOL: aws.Bool(false)}

	dynamoDBClient.On("UpdateItem", mock.Anything).Return(mockUpdateItemOutput, nil)

	result, err := table.UpdateOutput(mockUpdateItemAlertOutput)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.NotNil(t, err.(*genericapi.InternalError))
	dynamoDBClient.AssertExpectations(t)
}
