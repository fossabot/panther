package users

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/pkg/genericapi"
)

func (m *mockDynamoClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func TestDeleteDoesNotExist(t *testing.T) {
	mockClient := &mockDynamoClient{}
	returnErr := awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "", nil)
	mockClient.On("DeleteItem", mock.Anything).Return(
		(*dynamodb.DeleteItemOutput)(nil), returnErr)
	table := &Table{client: mockClient, Name: aws.String("table-name")}

	err := table.Delete(aws.String("id"))
	mockClient.AssertExpectations(t)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.DoesNotExistError{}, err)
}

func TestDeleteOtherError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	mockClient.On("DeleteItem", mock.Anything).Return(
		(*dynamodb.DeleteItemOutput)(nil), errors.New("service unavailable"))
	table := &Table{client: mockClient, Name: aws.String("table-name")}

	err := table.Delete(aws.String("id"))
	mockClient.AssertExpectations(t)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.AWSError{}, err)
}

func TestDelete(t *testing.T) {
	mockClient := &mockDynamoClient{}
	mockClient.On("DeleteItem", mock.Anything).Return(&dynamodb.DeleteItemOutput{}, nil)
	table := &Table{client: mockClient, Name: aws.String("table-name")}

	assert.NoError(t, table.Delete(aws.String("id")))
	mockClient.AssertExpectations(t)
}