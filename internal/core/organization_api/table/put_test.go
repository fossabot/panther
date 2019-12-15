package table

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

func (m *mockDynamoClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func TestPutItemError(t *testing.T) {
	mockClient := &mockDynamoClient{}
	mockClient.On("PutItem", mock.Anything).Return(
		(*dynamodb.PutItemOutput)(nil), errors.New("service unavailable"))
	table := &OrganizationsTable{client: mockClient, Name: aws.String("table-name")}

	err := table.Put(&models.Organization{})
	mockClient.AssertExpectations(t)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.AWSError{}, err)
}

func TestPutItem(t *testing.T) {
	mockClient := &mockDynamoClient{}
	mockClient.On("PutItem", mock.Anything).Return((*dynamodb.PutItemOutput)(nil), nil)
	table := &OrganizationsTable{client: mockClient, Name: aws.String("table-name")}

	assert.NoError(t, table.Put(&models.Organization{}))
	mockClient.AssertExpectations(t)
}
