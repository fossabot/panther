package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"

	"github.com/panther-labs/panther/api/lambda/users/models"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
)

var removeUserInput = &models.RemoveUserInput{
	ID:         aws.String("user123"),
	UserPoolID: aws.String("fakePoolId"),
}

func TestRemoveUserGetErr(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	mockGateway.On("GetUser", removeUserInput.ID, removeUserInput.UserPoolID).Return(&models.User{}, &genericapi.AWSError{})

	err := (API{}).RemoveUser(removeUserInput)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})

	mockGateway.AssertExpectations(t)
	m.AssertExpectations(t)
	mockGateway.AssertNotCalled(t, "DeleteUser")
	m.AssertNotCalled(t, "Delete")
}

func TestRemoveUserCognitoErr(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	mockGateway.On("GetUser", removeUserInput.ID, removeUserInput.UserPoolID).Return(&models.User{
		Email: aws.String("email@email.com"),
	}, nil)
	mockGateway.On("DeleteUser", removeUserInput.ID, removeUserInput.UserPoolID).Return(&genericapi.AWSError{})

	err := (API{}).RemoveUser(removeUserInput)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})

	mockGateway.AssertExpectations(t)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "Delete")
}

func TestRemoveUserDynamoErr(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	mockGateway.On("GetUser", removeUserInput.ID, removeUserInput.UserPoolID).Return(&models.User{
		Email: aws.String("email@email.com"),
	}, nil)
	mockGateway.On("DeleteUser", removeUserInput.ID, removeUserInput.UserPoolID).Return(nil)
	m.On("Delete", aws.String("email@email.com")).Return(&genericapi.AWSError{})

	err := (API{}).RemoveUser(removeUserInput)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})

	mockGateway.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestRemoveUserHandle(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	mockGateway.On("GetUser", removeUserInput.ID, removeUserInput.UserPoolID).Return(&models.User{
		Email: aws.String("email@email.com"),
	}, nil)
	mockGateway.On("DeleteUser", removeUserInput.ID, removeUserInput.UserPoolID).Return(nil)
	m.On("Delete", aws.String("email@email.com")).Return(nil)

	assert.NoError(t, (API{}).RemoveUser(removeUserInput))
	mockGateway.AssertExpectations(t)
	m.AssertExpectations(t)
}
