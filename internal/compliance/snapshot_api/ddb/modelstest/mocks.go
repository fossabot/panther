package modelstest

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

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/mock"
)

// MockDDBClient is used to stub out requests to DynamoDB for unit testing.
type MockDDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
	MockScanAttributes      []map[string]*dynamodb.AttributeValue
	MockItemAttributeOutput map[string]*dynamodb.AttributeValue
	MockQueryAttributes     []map[string]*dynamodb.AttributeValue
	TestErr                 bool
}

// DeleteItem is a mock method to remove an item from a dynamodb table.
func (client *MockDDBClient) DeleteItem(
	input *dynamodb.DeleteItemInput,
) (*dynamodb.DeleteItemOutput, error) {

	args := client.Called(input)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

// UpdateItem is a mock method to update an item from a dynamodb table.
func (client *MockDDBClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	args := client.Called(input)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

// Scan is a mock DynamoDB Scan request.
func (client *MockDDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	if client.TestErr {
		return nil, errors.New("fake dynamodb.Scan error")
	}
	return &dynamodb.ScanOutput{Items: client.MockScanAttributes}, nil
}

// Query is a mock DynamoDB Query request.
func (client *MockDDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if client.TestErr {
		return nil, errors.New("fake dynamodb.Query error")
	}
	return &dynamodb.QueryOutput{Items: client.MockQueryAttributes}, nil
}

// PutItem is a mock DynamoDB PutItem request.
func (client *MockDDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if client.TestErr {
		return nil, errors.New("fake dynamodb.PutItem error")
	}
	return &dynamodb.PutItemOutput{Attributes: client.MockItemAttributeOutput}, nil
}

func (client *MockDDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := client.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

// BatchWriteItem is a mock DynamoDB BatchWriteItem request.
func (client *MockDDBClient) BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	if client.TestErr {
		return nil, errors.New("fake dynamodb.BatchWriteItem error")
	}
	return &dynamodb.BatchWriteItemOutput{}, nil
}
