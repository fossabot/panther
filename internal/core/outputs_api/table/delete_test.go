package table

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
