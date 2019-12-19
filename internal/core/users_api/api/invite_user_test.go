package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var input = &models.InviteUserInput{
	GivenName:  aws.String("Joe"),
	Email:      aws.String("joe.blow@panther.io"),
	FamilyName: aws.String("Blow"),
	UserPoolID: aws.String("fakePoolId"),
	Role:       aws.String("Admin"),
}
var userID = aws.String("1234-5678-9012")

func TestInviteUserAddToOrgErr(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	m.On("Get", input.Email).Return(
		(*models.UserItem)(nil), &genericapi.AWSError{})

	// call the code we are testing
	result, err := (API{}).InviteUser(input)

	// assert that the expectations were met
	mockGateway.AssertExpectations(t)
	mockGateway.AssertNotCalled(t, "CreateUser")
	mockGateway.AssertNotCalled(t, "AddUserToGroup")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
}

func TestInviteUserAddToGroupErr(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	// setup gateway expectations
	m.On("Get", input.Email).Return((*models.UserItem)(nil), &genericapi.DoesNotExistError{})
	m.On("Put", &models.UserItem{
		ID: input.Email,
	}).Return(nil)
	mockGateway.On("CreateUser", &gateway.CreateUserInput{
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
		Email:      input.Email,
		UserPoolID: input.UserPoolID,
	}).Return(userID, nil)
	mockGateway.On("AddUserToGroup", userID, input.Role, input.UserPoolID).Return(&genericapi.AWSError{})

	// call the code we are testing
	result, err := (API{}).InviteUser(input)

	// assert that the expectations were met
	mockGateway.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
}

func TestInviteUserCreateErr(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	// setup gateway expectations
	m.On("Get", input.Email).Return((*models.UserItem)(nil), &genericapi.DoesNotExistError{})
	m.On("Put", &models.UserItem{
		ID: input.Email,
	}).Return(nil)
	mockGateway.On("CreateUser", &gateway.CreateUserInput{
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
		Email:      input.Email,
		UserPoolID: input.UserPoolID,
	}).Return(aws.String(""), &genericapi.AWSError{})
	m.On("Delete", input.Email).Return(nil)

	// call the code we are testing
	result, err := (API{}).InviteUser(input)

	// assert that the expectations were met
	mockGateway.AssertExpectations(t)
	mockGateway.AssertNotCalled(t, "AddUserToGroup")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
}

func TestInviteUserDeleteErr(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	// setup expectations
	m.On("Get", input.Email).Return((*models.UserItem)(nil), &genericapi.DoesNotExistError{})
	m.On("Put", &models.UserItem{
		ID: input.Email,
	}).Return(nil)
	mockGateway.On("CreateUser", &gateway.CreateUserInput{
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
		Email:      input.Email,
		UserPoolID: input.UserPoolID,
	}).Return(aws.String(""), &genericapi.AWSError{})
	m.On("Delete", input.Email).Return(&genericapi.AWSError{})

	// call the code we are testing
	result, err := (API{}).InviteUser(input)

	// assert that the expectations were met
	mockGateway.AssertExpectations(t)
	mockGateway.AssertNotCalled(t, "AddUserToGroup")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
}

func TestInviteUserHandle(t *testing.T) {
	// create an instance of our test objects
	mockGateway := &gateway.MockUserGateway{}
	m := &users.MockTable{}
	// replace the global variables with our mock objects
	userGateway = mockGateway
	userTable = m

	// setup gateway expectations
	m.On("Get", input.Email).Return((*models.UserItem)(nil), &genericapi.DoesNotExistError{})
	m.On("Put", &models.UserItem{
		ID: input.Email,
	}).Return(nil)
	mockGateway.On("CreateUser", &gateway.CreateUserInput{
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
		Email:      input.Email,
		UserPoolID: input.UserPoolID,
	}).Return(userID, nil)
	mockGateway.On("AddUserToGroup", userID, input.Role, input.UserPoolID).Return(nil)

	// call the code we are testing
	result, err := (API{}).InviteUser(input)

	// assert that the expectations were met
	mockGateway.AssertExpectations(t)
	assert.NotNil(t, result)
	assert.Equal(t, result, &models.InviteUserOutput{
		ID: userID,
	})
	assert.NoError(t, err)
}
