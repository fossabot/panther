package users

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/pkg/genericapi"
)

func (m *mockDynamoClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func TestGetItemAwsError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	mockClient.On("GetItem", mock.Anything).Return(
		(*dynamodb.GetItemOutput)(nil), errors.New("service unavailable"))
	table := &Table{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Get(aws.String("test-id"))
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.AWSError{}, err)
}

func TestGetItemUnmarshalError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	// output has wrong type for one of the fields
	output := &dynamodb.GetItemOutput{
		Item: DynamoItem{
			"id": {SS: aws.StringSlice([]string{"panther", "labs"})},
		},
	}
	mockClient.On("GetItem", mock.Anything).Return(output, nil)
	table := &Table{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Get(aws.String("test-id"))
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.InternalError{}, err)
}

func TestGetItemDoesNotExistError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	// output has wrong type for one of the fields
	output := &dynamodb.GetItemOutput{Item: DynamoItem{}}
	mockClient.On("GetItem", mock.Anything).Return(output, nil)
	table := &Table{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Get(aws.String("test-id"))
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.DoesNotExistError{}, err)
}

func TestGetItem(t *testing.T) {
	mockClient := &mockDynamoClient{}
	expectedInput := &dynamodb.GetItemInput{
		Key:       DynamoItem{"id": {S: aws.String("test-id")}},
		TableName: aws.String("test-table"),
	}
	// output has wrong type for one of the fields
	output := &dynamodb.GetItemOutput{Item: DynamoItem{"id": {S: aws.String("test-id")}}}
	mockClient.On("GetItem", expectedInput).Return(output, nil)
	table := &Table{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Get(aws.String("test-id"))
	mockClient.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, aws.String("test-id"), result.ID)
}
