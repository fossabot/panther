package modelstest

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
func (client *MockDDBClient) UpdateItem(
	input *dynamodb.UpdateItemInput,
) (*dynamodb.UpdateItemOutput, error) {

	args := client.Called(input)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

// Scan is a mock DynamoDB Scan request.
func (client *MockDDBClient) Scan(
	input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {

	if client.TestErr {
		return nil, errors.New("fake dynamodb.Scan error")
	}

	return &dynamodb.ScanOutput{Items: client.MockScanAttributes}, nil
}

// Query is a mock DynamoDB Query request.
func (client *MockDDBClient) Query(
	input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {

	if client.TestErr {
		return nil, errors.New("fake dynamodb.Query error")
	}

	return &dynamodb.QueryOutput{Items: client.MockQueryAttributes}, nil
}

// PutItem is a mock DynamoDB PutItem request.
func (client *MockDDBClient) PutItem(
	input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {

	if client.TestErr {
		return nil, errors.New("fake dynamodb.PutItem error")
	}

	return &dynamodb.PutItemOutput{Attributes: client.MockItemAttributeOutput}, nil
}

// GetItem is a mock DynamoDB GetItem request.
func (client *MockDDBClient) GetItem(
	input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	if client.TestErr {
		return nil, errors.New("fake dynamodb.GetItem error")
	}

	return &dynamodb.GetItemOutput{Item: client.MockItemAttributeOutput}, nil
}

// BatchWriteItem is a mock DynamoDB BatchWriteItem request.
func (client *MockDDBClient) BatchWriteItem(
	input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {

	if client.TestErr {
		return nil, errors.New("fake dynamodb.BatchWriteItem error")
	}

	return &dynamodb.BatchWriteItemOutput{}, nil
}