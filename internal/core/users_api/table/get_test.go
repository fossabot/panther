package users

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
