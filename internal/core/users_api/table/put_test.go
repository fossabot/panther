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

	"github.com/panther-labs/panther/api/lambda/users/models"
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
	table := &Table{client: mockClient, Name: aws.String("table-name")}

	err := table.Put(&models.UserItem{
		ID: aws.String("test-id"),
	})
	mockClient.AssertExpectations(t)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.AWSError{}, err)
}

func TestPutItem(t *testing.T) {
	mockClient := &mockDynamoClient{}
	mockClient.On("PutItem", mock.Anything).Return((*dynamodb.PutItemOutput)(nil), nil)
	table := &Table{client: mockClient, Name: aws.String("table-name")}

	assert.NoError(t, table.Put(&models.UserItem{
		ID: aws.String("test-id"),
	}))
	mockClient.AssertExpectations(t)
}
