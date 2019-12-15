package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/api/lambda/users/models"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
)

var testAddUserToOrgInput = &models.AddUserToOrganizationInput{
	Email: aws.String("email@emails.com"),
}

func TestAddUserToOrgAlreadyExists(t *testing.T) {
	// create an instance of our test objects
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userTable = m

	m.On("Get", testAddUserToOrgInput.Email).Return(
		&models.UserItem{}, nil)

	result, err := (API{}).AddUserToOrganization(testAddUserToOrgInput)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.AlreadyExistsError{}, err)
	assert.Nil(t, result)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "Put")
}

func TestAddUserToOrgGetDynamoError(t *testing.T) {
	// create an instance of our test objects
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userTable = m

	m.On("Get", testAddUserToOrgInput.Email).Return(
		(*models.UserItem)(nil), &genericapi.LambdaError{})

	result, err := (API{}).AddUserToOrganization(testAddUserToOrgInput)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.LambdaError{}, err)
	assert.Nil(t, result)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "Put")
}

func TestAddUserToOrgPutError(t *testing.T) {
	// create an instance of our test objects
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userTable = m

	m.On("Get", testAddUserToOrgInput.Email).Return(
		(*models.UserItem)(nil), &genericapi.DoesNotExistError{})
	m.On("Put", &models.UserItem{
		ID: testAddUserToOrgInput.Email,
	}).Return(&genericapi.AWSError{})

	result, err := (API{}).AddUserToOrganization(testAddUserToOrgInput)
	assert.Error(t, err)
	assert.IsType(t, &genericapi.AWSError{}, err)
	assert.Nil(t, result)
	m.AssertExpectations(t)
}

func TestAddUserToOrgHandle(t *testing.T) {
	// create an instance of our test objects
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userTable = m

	m.On("Get", testAddUserToOrgInput.Email).Return(
		(*models.UserItem)(nil), &genericapi.DoesNotExistError{})
	m.On("Put", &models.UserItem{
		ID: testAddUserToOrgInput.Email,
	}).Return(nil)

	result, err := (API{}).AddUserToOrganization(testAddUserToOrgInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result, &models.AddUserToOrganizationOutput{
		Email: testAddUserToOrgInput.Email,
	})
	m.AssertExpectations(t)
}
