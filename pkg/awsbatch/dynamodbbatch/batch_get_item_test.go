package dynamodbbatch

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (m *mockDynamo) BatchGetItemPages(
	input *dynamodb.BatchGetItemInput,
	pageFunc func(*dynamodb.BatchGetItemOutput, bool) bool,
) error {
	m.callCount++

	if m.err != nil {
		return m.err
	}

	// Put all items in a single response page
	page := &dynamodb.BatchGetItemOutput{
		Responses: map[string][]map[string]*dynamodb.AttributeValue{
			mockTableName: input.RequestItems[mockTableName].Keys,
		},
	}
	pageFunc(page, true)
	return nil
}

func mockGetInput() *dynamodb.BatchGetItemInput {
	return &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			mockTableName: {
				Keys: []map[string]*dynamodb.AttributeValue{
					{"eventID": &dynamodb.AttributeValue{S: aws.String("id-1")}},
					{"eventID": &dynamodb.AttributeValue{S: aws.String("id-2")}},
				},
			},
		},
	}
}

func TestBatchGetItemCount(t *testing.T) {
	items := map[string]*dynamodb.KeysAndAttributes{
		"table1": {Keys: make([]map[string]*dynamodb.AttributeValue, 2)},
		"table2": {Keys: make([]map[string]*dynamodb.AttributeValue, 3)},
		"table3": {Keys: make([]map[string]*dynamodb.AttributeValue, 5)},
	}
	assert.Equal(t, 10, getItemCount(items))
}

func TestBatchGetItem(t *testing.T) {
	client := &mockDynamo{}
	result, err := BatchGetItem(client, mockGetInput())
	require.Nil(t, err)
	assert.Equal(t, 1, client.callCount)
	assert.Equal(t, 2, len(result.Responses[mockTableName]))
	assert.Equal(t, "id-1", *result.Responses[mockTableName][0]["eventID"].S)
	assert.Equal(t, "id-2", *result.Responses[mockTableName][1]["eventID"].S)
}

// An error is returned
func TestBatchGetItemError(t *testing.T) {
	client := &mockDynamo{err: errors.New("internal service error")}
	result, err := BatchGetItem(client, mockGetInput())
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

// A large number of records are broken into multiple requests
func TestBatchGetItemPaging(t *testing.T) {
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			mockTableName: {
				Keys: make([]map[string]*dynamodb.AttributeValue, maxBatchGetItems*2+1),
			},
		},
	}
	client := &mockDynamo{}
	result, err := BatchGetItem(client, input)
	require.Nil(t, err)
	assert.Equal(t, 3, client.callCount)
	assert.Equal(t, maxBatchGetItems*2+1, len(result.Responses[mockTableName]))
}
