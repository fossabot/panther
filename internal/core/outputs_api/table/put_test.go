package table

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var mockOutputID = aws.String("outputID")

type mockPutClient struct {
	dynamodbiface.DynamoDBAPI
	conditionalErr bool
	serviceErr     bool
}

func (m *mockPutClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if m.conditionalErr && input.ConditionExpression != nil {
		return nil, awserr.New(
			dynamodb.ErrCodeConditionalCheckFailedException, "attribute does not exist", nil)
	}
	if m.serviceErr {
		return nil, awserr.New(
			dynamodb.ErrCodeResourceNotFoundException, "table does not exist", nil)
	}
	return &dynamodb.PutItemOutput{}, nil
}

func TestPutOutputDoesNotExist(t *testing.T) {
	table := &OutputsTable{client: &mockPutClient{conditionalErr: true}}
	err := table.PutOutput(&models.AlertOutputItem{OutputID: mockOutputID})
	assert.NotNil(t, err.(*genericapi.DoesNotExistError))
}

func TestPutOutputServiceError(t *testing.T) {
	table := &OutputsTable{client: &mockPutClient{serviceErr: true}}
	err := table.PutOutput(&models.AlertOutputItem{OutputID: mockOutputID})
	assert.NotNil(t, err.(*genericapi.AWSError))
}

func TestPutOutput(t *testing.T) {
	table := &OutputsTable{client: &mockPutClient{}}
	assert.Nil(t, table.PutOutput(&models.AlertOutputItem{OutputID: mockOutputID}))
}
