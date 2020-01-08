package api

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

var mockGetDefaultsInput = &models.GetDefaultOutputsInput{}

func TestGetDefaultOutputs(t *testing.T) {
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	items := []*models.DefaultOutputsItem{
		{
			Severity:  aws.String("INFO"),
			OutputIDs: []*string{aws.String("outputId1")},
		},
		{
			Severity:  aws.String("WARN"),
			OutputIDs: []*string{aws.String("outputId2")},
		},
	}

	expectedResult := &models.GetDefaultOutputsOutput{
		Defaults: []*models.DefaultOutputs{
			{
				Severity:  aws.String("INFO"),
				OutputIDs: []*string{aws.String("outputId1")},
			},
			{
				Severity:  aws.String("WARN"),
				OutputIDs: []*string{aws.String("outputId2")},
			},
		},
	}
	mockDefaultsTable.On("GetDefaults").Return(items, nil)

	result, err := (API{}).GetDefaultOutputs(mockGetDefaultsInput)

	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockDefaultsTable.AssertExpectations(t)
}

func TestGetDefaultOutputsEmpty(t *testing.T) {
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	items := []*models.DefaultOutputsItem{
		{
			Severity: aws.String("HIGH"),
			// DDB returns the below structure instead of populated empty slice
			OutputIDs: nil,
		},
	}
	// Verify that the result has the an empty slice in the Default field instead of nil
	expectedResult := &models.GetDefaultOutputsOutput{Defaults: []*models.DefaultOutputs{}}

	mockDefaultsTable.On("GetDefaults").Return(items, nil)

	result, err := (API{}).GetDefaultOutputs(mockGetDefaultsInput)

	require.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockDefaultsTable.AssertExpectations(t)
}

func TestGetDefaultOutputsTableError(t *testing.T) {
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable
	mockDefaultsTable.On("GetDefaults").Return(([]*models.DefaultOutputsItem)(nil), errors.New("error"))

	result, err := (API{}).GetDefaultOutputs(mockGetDefaultsInput)

	require.Error(t, err)
	assert.Nil(t, result)
	mockDefaultsTable.AssertExpectations(t)
}
