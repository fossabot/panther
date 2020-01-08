package gateway

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

// MockUserGateway is a mocked object that implements the API interface
// It describes an object that the apis rely on.
type MockUserGateway struct {
	API
	mock.Mock
}

// The following methods implement the API interface
// and just record the activity, and returns what the Mock object tells it to.

// AddUserToGroup mocks AddUserToGroup for testing
func (m *MockUserGateway) AddUserToGroup(id *string, groupName *string, userPoolID *string) error {
	args := m.Called(id, groupName, userPoolID)
	return args.Error(0)
}

// CreateUser mocks CreateUser for testing
func (m *MockUserGateway) CreateUser(input *CreateUserInput) (*string, error) {
	args := m.Called(input)
	return args.Get(0).(*string), args.Error(1)
}

// CreateUserPool mocks CreateUserPool for testing
func (m *MockUserGateway) CreateUserPool(displayName *string) (*UserPool, error) {
	args := m.Called(displayName)
	return args.Get(0).(*UserPool), args.Error(1)
}

// CreateUserPoolGroups mocks CreateUserPoolGroups for testing
func (m *MockUserGateway) CreateUserPoolGroups(userPoolID *string) error {
	args := m.Called(userPoolID)
	return args.Error(0)
}

// DeleteUser mocks DeleteUser for testing
func (m *MockUserGateway) DeleteUser(id *string, userPoolID *string) error {
	args := m.Called(id, userPoolID)
	return args.Error(0)
}

// GetUser mocks GetUser for testing
func (m *MockUserGateway) GetUser(id *string, userPoolID *string) (*models.User, error) {
	args := m.Called(id, userPoolID)
	return args.Get(0).(*models.User), args.Error(1)
}

// ListGroups mocks ListGroups for testing
func (m *MockUserGateway) ListGroups(userPoolID *string) ([]*models.Group, error) {
	args := m.Called(userPoolID)
	return args.Get(0).([]*models.Group), args.Error(1)
}

// ListGroupsForUser mocks ListGroupsForUser for testing
func (m *MockUserGateway) ListGroupsForUser(id *string, userPoolID *string) ([]*models.Group, error) {
	args := m.Called(id, userPoolID)
	return args.Get(0).([]*models.Group), args.Error(1)
}

// ListUsers mocks ListUsers for testing
func (m *MockUserGateway) ListUsers(limit *int64, paginationToken *string, userPoolID *string) (*ListUsersOutput, error) {
	args := m.Called(limit, paginationToken, userPoolID)
	return args.Get(0).(*ListUsersOutput), args.Error(1)
}

// RemoveUserFromGroup mocks RemoveUserFromGroup for testing
func (m *MockUserGateway) RemoveUserFromGroup(id *string, groupName *string, userPoolID *string) error {
	args := m.Called(id, groupName, userPoolID)
	return args.Error(0)
}

// ResetUserPassword mocks ResetUserPassword for testing
func (m *MockUserGateway) ResetUserPassword(id *string, userPoolID *string) error {
	args := m.Called(id, userPoolID)
	return args.Error(0)
}

// UpdateUser mocks UpdateUser for testing
func (m *MockUserGateway) UpdateUser(input *UpdateUserInput) error {
	args := m.Called(input)
	return args.Error(0)
}
