package api

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb/modelstest"
	"github.com/panther-labs/panther/pkg/genericapi"
)

func TestDeleteIntegrationItem(t *testing.T) {
	mockClient := &modelstest.MockDDBClient{}
	db = &ddb.DDB{Client: mockClient, TableName: "test"}

	resp := &dynamodb.DeleteItemOutput{}
	mockClient.On("DeleteItem", mock.Anything).Return(resp, nil)

	result := apiTest.DeleteIntegration(&models.DeleteIntegrationInput{
		IntegrationID: aws.String(testIntegrationID),
	})

	assert.NoError(t, result)
	mockClient.AssertExpectations(t)
}

func TestDeleteIntegrationItemError(t *testing.T) {
	mockClient := &modelstest.MockDDBClient{}
	db = &ddb.DDB{Client: mockClient, TableName: "test"}

	mockErr := awserr.New(
		"ErrCodeInternalServerError",
		"An error occurred on the server side.",
		errors.New("fake error"),
	)
	resp := &dynamodb.DeleteItemOutput{}
	mockClient.On("DeleteItem", mock.Anything).Return(resp, mockErr)

	result := apiTest.DeleteIntegration(&models.DeleteIntegrationInput{
		IntegrationID: aws.String(testIntegrationID),
	})

	assert.Error(t, result)
	mockClient.AssertExpectations(t)
}

func TestDeleteIntegrationItemDoesNotExist(t *testing.T) {
	mockClient := &modelstest.MockDDBClient{}
	db = &ddb.DDB{Client: mockClient, TableName: "test"}

	mockErr := awserr.New(
		"ConditionalCheckFailedException",
		"A condition specified in the operation could not be evaluated.",
		errors.New("fake error"),
	)
	resp := &dynamodb.DeleteItemOutput{}
	mockClient.On("DeleteItem", mock.Anything).Return(resp, mockErr)

	result := apiTest.DeleteIntegration(&models.DeleteIntegrationInput{
		IntegrationID: aws.String(testIntegrationID),
	})

	assert.Error(t, result)
	assert.IsType(t, &genericapi.DoesNotExistError{}, result)
	mockClient.AssertExpectations(t)
}
