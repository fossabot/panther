package table

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

var mockDefaultInputsItem = &models.DefaultOutputsItem{
	Severity:  aws.String("INFO"),
	OutputIDs: []*string{aws.String("outputId")},
}

func TestPutDefaults(t *testing.T) {
	mockClient := &mockDynamoDB{}
	table := &DefaultsTable{client: mockClient, Name: aws.String("defaultsTable")}

	expectedPutItem := &dynamodb.PutItemInput{
		TableName: aws.String("defaultsTable"),
		Item: map[string]*dynamodb.AttributeValue{
			"severity": {
				S: aws.String("INFO"),
			},
			"outputIds": {
				SS: aws.StringSlice([]string{"outputId"}),
			},
		},
	}

	mockClient.On("PutItem", expectedPutItem).Return((*dynamodb.PutItemOutput)(nil), nil)

	require.NoError(t, table.PutDefaults(mockDefaultInputsItem))
	mockClient.AssertExpectations(t)
}

func TestPutDefaultsClientError(t *testing.T) {
	mockClient := &mockDynamoDB{}
	table := &DefaultsTable{client: mockClient, Name: aws.String("defaultsTable")}
	mockClient.On("PutItem", mock.Anything).Return((*dynamodb.PutItemOutput)(nil), errors.New("testing"))

	require.Error(t, table.PutDefaults(mockDefaultInputsItem))
	mockClient.AssertExpectations(t)
}
