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
