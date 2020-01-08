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
