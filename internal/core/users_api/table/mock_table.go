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
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// MockTable is a mocked object that implements the User Table API interface
type MockTable struct {
	API
	mock.Mock
}

// AddUserToOrganization mocks AddUserToOrganization for testing
func (m *MockTable) AddUserToOrganization(userItem *models.UserItem) error {
	args := m.Called(userItem)
	return args.Error(0)
}

// Delete mocks Delete for testing
func (m *MockTable) Delete(id *string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Get mocks Get for testing
func (m *MockTable) Get(id *string) (*models.UserItem, error) {
	args := m.Called(id)
	return args.Get(0).(*models.UserItem), args.Error(1)
}

// Put mocks Put for testing
func (m *MockTable) Put(userItem *models.UserItem) error {
	args := m.Called(userItem)
	return args.Error(0)
}
