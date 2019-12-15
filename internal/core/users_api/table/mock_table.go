package users

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
