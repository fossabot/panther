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

	"github.com/panther-labs/panther/pkg/genericapi"
)

var mockDeleteItemOutput = &dynamodb.DeleteItemOutput{}
var deleteOutputID = aws.String("outputId")

func TestDeleteOutput(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}

	expectedCondition := expression.Name("outputId").Equal(expression.Value(aws.String("outputId")))

	expectedConditionExpression, _ := expression.NewBuilder().WithCondition(expectedCondition).Build()

	expectedDeleteItemInput := &dynamodb.DeleteItemInput{
		Key: DynamoItem{
			"outputId": {S: aws.String("outputId")},
		},
		TableName:                 aws.String("TableName"),
		ConditionExpression:       expectedConditionExpression.Condition(),
		ExpressionAttributeNames:  expectedConditionExpression.Names(),
		ExpressionAttributeValues: expectedConditionExpression.Values(),
	}

	dynamoDBClient.On("DeleteItem", expectedDeleteItemInput).Return(mockDeleteItemOutput, nil)

	assert.NoError(t, table.DeleteOutput(deleteOutputID))
	dynamoDBClient.AssertExpectations(t)
}

func TestDeleteOutputDoesNotExist(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}

	dynamoDBClient.On("DeleteItem", mock.Anything).Return(
		mockDeleteItemOutput,
		awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "attribute does not exist", nil))

	result := table.DeleteOutput(deleteOutputID)
	assert.Error(t, result)
	assert.NotNil(t, result.(*genericapi.DoesNotExistError))
	dynamoDBClient.AssertExpectations(t)
}

func TestDeleteOutputServiceError(t *testing.T) {
	dynamoDBClient := &mockDynamoDB{}
	table := &OutputsTable{client: dynamoDBClient, Name: aws.String("TableName")}

	dynamoDBClient.On("DeleteItem", mock.Anything).Return(
		mockDeleteItemOutput,
		awserr.New(dynamodb.ErrCodeResourceNotFoundException, "table does not exist", nil))

	result := table.DeleteOutput(deleteOutputID)
	assert.Error(t, result)
	assert.NotNil(t, result.(*genericapi.AWSError))
	dynamoDBClient.AssertExpectations(t)
}
