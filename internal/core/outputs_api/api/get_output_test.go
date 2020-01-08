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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

func TestGetOutput(t *testing.T) {
	mockOutputsTable := &mockOutputTable{}
	outputsTable = mockOutputsTable
	mockEncryptionKey := &mockEncryptionKey{}
	encryptionKey = mockEncryptionKey
	mockOutputVerification := &mockOutputVerification{}
	outputVerification = mockOutputVerification
	mockDefaultsTable := &mockDefaultsTable{}
	defaultsTable = mockDefaultsTable

	mockGetOutputInput := &models.GetOutputInput{
		OutputID: aws.String("outputId"),
	}

	alertOutputItem := &models.AlertOutputItem{
		OutputID:           aws.String("outputId"),
		DisplayName:        aws.String("displayName"),
		CreatedBy:          aws.String("createdBy"),
		CreationTime:       aws.String("creationTime"),
		LastModifiedBy:     aws.String("lastModifiedBy"),
		LastModifiedTime:   aws.String("lastModifiedTime"),
		OutputType:         aws.String("slack"),
		EncryptedConfig:    make([]byte, 1),
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
	}
	mockEncryptionKey.On("DecryptConfig", make([]byte, 1), mock.Anything).Return(nil)
	mockOutputsTable.On("GetOutput", aws.String("outputId")).Return(alertOutputItem, nil)
	mockDefaultsTable.On("GetDefaults", mock.Anything).Return([]*models.DefaultOutputsItem{}, nil)

	expectedAlertOutput := &models.AlertOutput{
		OutputID:           aws.String("outputId"),
		OutputType:         aws.String("slack"),
		CreatedBy:          aws.String("createdBy"),
		CreationTime:       aws.String("creationTime"),
		DisplayName:        aws.String("displayName"),
		LastModifiedBy:     aws.String("lastModifiedBy"),
		LastModifiedTime:   aws.String("lastModifiedTime"),
		OutputConfig:       &models.OutputConfig{},
		VerificationStatus: aws.String(models.VerificationStatusSuccess),
		DefaultForSeverity: []*string{},
	}

	result, err := (API{}).GetOutput(mockGetOutputInput)

	require.NoError(t, err)
	assert.Equal(t, expectedAlertOutput, result)
	mockOutputsTable.AssertExpectations(t)
	mockEncryptionKey.AssertExpectations(t)
	mockDefaultsTable.AssertExpectations(t)
}
