package table

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
