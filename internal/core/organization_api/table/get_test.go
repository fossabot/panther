package table

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
	table := &OrganizationsTable{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Get()
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
			"awsConfig": {SS: aws.StringSlice([]string{"panther", "labs"})},
		},
	}
	mockClient.On("GetItem", mock.Anything).Return(output, nil)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Get()
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.InternalError{}, err)
}

func TestGetItemDoesNotExistError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	output := &dynamodb.GetItemOutput{Item: DynamoItem{}}
	mockClient.On("GetItem", mock.Anything).Return(output, nil)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("test-table")}

	result, err := table.Get()
	mockClient.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.DoesNotExistError{}, err)
}

func TestGetItem(t *testing.T) {
	mockClient := &mockDynamoClient{}
	expectedInput := &dynamodb.GetItemInput{
		Key:       DynamoItem{"id": {S: aws.String("1")}},
		TableName: aws.String("test-table"),
	}
	output := &dynamodb.GetItemOutput{Item: DynamoItem{
		"id":        {S: aws.String("1")},
		"awsConfig": dynamoAwsConfig,
	}}
	mockClient.On("GetItem", expectedInput).Return(output, nil)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("test-table")}

	_, err := table.Get()
	mockClient.AssertExpectations(t)
	require.NoError(t, err)
}