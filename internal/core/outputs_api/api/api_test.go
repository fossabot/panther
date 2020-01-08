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
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/internal/core/outputs_api/encryption"
	"github.com/panther-labs/panther/internal/core/outputs_api/table"
	"github.com/panther-labs/panther/internal/core/outputs_api/verification"
)

type mockOutputTable struct {
	table.OutputsTable
	mock.Mock
}

func (m *mockOutputTable) GetOutput(outputID *string) (*models.AlertOutputItem, error) {
	args := m.Called(outputID)
	return args.Get(0).(*models.AlertOutputItem), args.Error(1)
}

func (m *mockOutputTable) DeleteOutput(outputID *string) error {
	args := m.Called(outputID)
	return args.Error(0)
}

func (m *mockOutputTable) GetOutputs() ([]*models.AlertOutputItem, error) {
	args := m.Called()
	return args.Get(0).([]*models.AlertOutputItem), args.Error(1)
}

func (m *mockOutputTable) UpdateOutput(input *models.AlertOutputItem) (*models.AlertOutputItem, error) {
	args := m.Called(input)
	return args.Get(0).(*models.AlertOutputItem), args.Error(1)
}

func (m *mockOutputTable) GetOutputByName(displayName *string) (*models.AlertOutputItem, error) {
	args := m.Called(displayName)
	alertOutputItem := args.Get(0)
	if alertOutputItem == nil {
		return nil, args.Error(1)
	}
	return alertOutputItem.(*models.AlertOutputItem), args.Error(1)
}

func (m *mockOutputTable) PutOutput(output *models.AlertOutputItem) error {
	args := m.Called(output)
	return args.Error(0)
}

type mockDefaultsTable struct {
	table.DefaultsTable
	mock.Mock
}

func (m *mockDefaultsTable) GetDefaults() ([]*models.DefaultOutputsItem, error) {
	args := m.Called()
	return args.Get(0).([]*models.DefaultOutputsItem), args.Error(1)
}

func (m *mockDefaultsTable) GetDefault(severity *string) (*models.DefaultOutputsItem, error) {
	args := m.Called()
	return args.Get(0).(*models.DefaultOutputsItem), args.Error(1)
}

func (m *mockDefaultsTable) PutDefaults(item *models.DefaultOutputsItem) error {
	args := m.Called(item)
	return args.Error(0)
}

type mockEncryptionKey struct {
	encryption.Key
	mock.Mock
}

func (m *mockEncryptionKey) DecryptConfig(ciphertext []byte, config interface{}) error {
	args := m.Called(ciphertext, config)
	return args.Error(0)
}

func (m *mockEncryptionKey) EncryptConfig(config interface{}) ([]byte, error) {
	args := m.Called(config)
	return args.Get(0).([]byte), args.Error(1)
}

type mockOutputVerification struct {
	verification.OutputVerificationAPI
	mock.Mock
}

func (m *mockOutputVerification) GetVerificationStatus(output *models.AlertOutput) (*string, error) {
	args := m.Called(output)
	return args.Get(0).(*string), args.Error(1)
}

func (m *mockOutputVerification) VerifyOutput(output *models.AlertOutput) (*models.AlertOutput, error) {
	args := m.Called(output)
	return args.Get(0).(*models.AlertOutput), args.Error(1)
}
