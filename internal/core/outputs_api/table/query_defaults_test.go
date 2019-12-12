package table

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var mockQueryOutput = &dynamodb.QueryOutput{
	Items: []map[string]*dynamodb.AttributeValue{
		{
			"severity": {
				S: aws.String("INFO"),
			},
			"outputIds": {
				L: []*dynamodb.AttributeValue{{S: aws.String("outputId")}},
			},
		},
	},
}

var mockScanOutput = &dynamodb.ScanOutput{
	Items: []map[string]*dynamodb.AttributeValue{
		{
			"severity": {
				S: aws.String("INFO"),
			},
			"outputIds": {
				L: []*dynamodb.AttributeValue{{S: aws.String("outputId")}},
			},
		},
	},
}

func TestGetDefaults(t *testing.T) {
	mockClient := &mockDynamoDB{}
	table := &DefaultsTable{client: mockClient, Name: aws.String("defaultsTable")}

	expectedResult := []*models.DefaultOutputsItem{
		{
			Severity:  aws.String("INFO"),
			OutputIDs: []*string{aws.String("outputId")},
		},
	}
	mockClient.On("ScanPages", mock.Anything, mock.AnythingOfType("func(*dynamodb.ScanOutput, bool) bool")).Return(nil)

	result, err := table.GetDefaults()

	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetDefaultsClientError(t *testing.T) {
	mockClient := &mockDynamoDB{}
	table := &DefaultsTable{client: mockClient, Name: aws.String("defaultsTable")}

	mockClient.On("ScanPages", mock.Anything, mock.AnythingOfType("func(*dynamodb.ScanOutput, bool) bool")).Return(errors.New("error" +
		""))

	result, err := table.GetDefaults()
	require.Error(t, err)
	assert.IsType(t, &genericapi.AWSError{}, err)
	assert.Nil(t, result)
}
