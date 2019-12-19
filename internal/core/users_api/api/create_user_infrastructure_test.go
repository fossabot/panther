package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var testInput = &models.CreateUserInfrastructureInput{
	DisplayName: aws.String("CompanyCompany"),
	GivenName:   aws.String("Joe"),
	Email:       aws.String("joe.blow@panther.io"),
	FamilyName:  aws.String("Blow"),
}

var testCreateUserInput = &gateway.CreateUserInput{
	GivenName:  testInput.GivenName,
	FamilyName: testInput.FamilyName,
	Email:      testInput.Email,
	UserPoolID: aws.String("us-west-2_ZlG7Ldp1K"),
}

func TestCreateUserPoolSuccess(t *testing.T) {
	testUserPoolID := aws.String("us-west-2_ZlG7Ldp1K")
	testAppClientID := aws.String("abcdefghijklmnopq")
	testIdentityPoolID := aws.String("us-west-2:abcdefghijklmnopq")
	testUserPool := &gateway.UserPool{
		UserPoolID:     testUserPoolID,
		AppClientID:    testAppClientID,
		IdentityPoolID: testIdentityPoolID,
	}
	testUserID := aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766")
	testGroup := aws.String("Admins")

	// create an instance of our test object
	mockGateway := &gateway.MockUserGateway{}
	// replace the global variable with our mock object
	userGateway = mockGateway

	// setup expectations
	mockGateway.On("CreateUserPool", testInput.DisplayName).Return(testUserPool, nil)
	mockGateway.On("CreateUserPoolGroups", testUserPoolID).Return(nil)
	mockGateway.On("CreateUser", testCreateUserInput).Return(testUserID, nil)
	mockGateway.On("AddUserToGroup", testUserID, testGroup, testUserPoolID).Return(nil)
	mockGateway.On("GetUser", testUserID, testUserPoolID).Return(&models.User{ID: testUserID}, nil)

	// call the code we are testing
	result, err := (API{}).CreateUserInfrastructure(testInput)

	// assert that the expectations were met
	mockGateway.AssertExpectations(t)
	assert.NotNil(t, result)
	assert.Equal(t, result.UserPoolID, testUserPoolID)
	assert.Equal(t, result.User.ID, testUserID)
	assert.Equal(t, result.AppClientID, testAppClientID)
	assert.NoError(t, err)
}

func TestCreatePoolFailure(t *testing.T) {
	testUserPoolID := aws.String("us-west-2_ZlG7Ldp1K")
	testAppClientID := aws.String("abcdefghijklmnopq")
	testIdentityPoolID := aws.String("us-west-2:abcdefghijklmnopq")
	testUserPool := &gateway.UserPool{
		UserPoolID:     testUserPoolID,
		AppClientID:    testAppClientID,
		IdentityPoolID: testIdentityPoolID,
	}
	// create an instance of our test object
	mockGateway := &gateway.MockUserGateway{}
	// replace the global variable with our mock object
	userGateway = mockGateway

	// setup expectations
	mockGateway.On("CreateUserPool", testInput.DisplayName).Return(testUserPool, &genericapi.AWSError{})

	// call the code we are testing
	result, err := (API{}).CreateUserInfrastructure(testInput)

	// assert that the expectations were met
	mockGateway.AssertNotCalled(t, "CreateUserPoolGroups", mock.Anything)
	mockGateway.AssertNotCalled(t, "CreateUser", mock.Anything)
	mockGateway.AssertNotCalled(t, "AddUserToGroup", mock.Anything)
	mockGateway.AssertNotCalled(t, "AddUserPoolToAppSync", mock.Anything)
	mockGateway.AssertNotCalled(t, "GetUser", mock.Anything)
	mockGateway.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestCreateGroupsFailure(t *testing.T) {
	testUserPoolID := aws.String("us-west-2_ZlG7Ldp1K")
	testAppClientID := aws.String("abcdefghijklmnopq")
	testIdentityPoolID := aws.String("us-west-2:abcdefghijklmnopq")
	testUserPool := &gateway.UserPool{
		UserPoolID:     testUserPoolID,
		AppClientID:    testAppClientID,
		IdentityPoolID: testIdentityPoolID,
	}
	// create an instance of our test object
	mockGateway := &gateway.MockUserGateway{}
	// replace the global variable with our mock object
	userGateway = mockGateway

	// setup expectations
	mockGateway.On("CreateUserPool", testInput.DisplayName).Return(testUserPool, nil)
	mockGateway.On("CreateUserPoolGroups", testUserPoolID).Return(&genericapi.AWSError{})

	// call the code we are testing
	result, err := (API{}).CreateUserInfrastructure(testInput)

	// assert that the expectations were met
	mockGateway.AssertNotCalled(t, "CreateUser", mock.Anything)
	mockGateway.AssertNotCalled(t, "AddUserToGroup", mock.Anything)
	mockGateway.AssertNotCalled(t, "AddUserPoolToAppSync", mock.Anything)
	mockGateway.AssertNotCalled(t, "GetUser", mock.Anything)
	mockGateway.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestCreateUserFailure(t *testing.T) {
	testUserPoolID := aws.String("us-west-2_ZlG7Ldp1K")
	testAppClientID := aws.String("abcdefghijklmnopq")
	testIdentityPoolID := aws.String("us-west-2:abcdefghijklmnopq")
	testUserPool := &gateway.UserPool{
		UserPoolID:     testUserPoolID,
		AppClientID:    testAppClientID,
		IdentityPoolID: testIdentityPoolID,
	}
	testUserID := aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766")

	// create an instance of our test object
	mockGateway := &gateway.MockUserGateway{}
	// replace the global variable with our mock object
	userGateway = mockGateway

	// setup expectations
	mockGateway.On("CreateUserPool", testInput.DisplayName).Return(testUserPool, nil)
	mockGateway.On("CreateUserPoolGroups", testUserPoolID).Return(nil)
	mockGateway.On("CreateUser", testCreateUserInput).Return(testUserID, &genericapi.AWSError{})
	// call the code we are testing
	result, err := (API{}).CreateUserInfrastructure(testInput)

	// assert that the expectations were met
	mockGateway.AssertNotCalled(t, "AddUserToGroup", mock.Anything)
	mockGateway.AssertNotCalled(t, "AddUserPoolToAppSync", mock.Anything)
	mockGateway.AssertNotCalled(t, "GetUser", mock.Anything)
	mockGateway.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestAddUserToGroupFailure(t *testing.T) {
	testUserPoolID := aws.String("us-west-2_ZlG7Ldp1K")
	testAppClientID := aws.String("abcdefghijklmnopq")
	testIdentityPoolID := aws.String("us-west-2:abcdefghijklmnopq")
	testUserPool := &gateway.UserPool{
		UserPoolID:     testUserPoolID,
		AppClientID:    testAppClientID,
		IdentityPoolID: testIdentityPoolID,
	}
	testUserID := aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766")
	testGroup := aws.String("Admins")

	// create an instance of our test object
	mockGateway := &gateway.MockUserGateway{}
	// replace the global variable with our mock object
	userGateway = mockGateway

	// setup expectations
	mockGateway.On("CreateUserPool", testInput.DisplayName).Return(testUserPool, nil)
	mockGateway.On("CreateUserPoolGroups", testUserPoolID).Return(nil)
	mockGateway.On("CreateUser", testCreateUserInput).Return(testUserID, nil)
	mockGateway.On("AddUserToGroup", testUserID, testGroup, testUserPoolID).Return(&genericapi.AWSError{})
	// call the code we are testing
	result, err := (API{}).CreateUserInfrastructure(testInput)

	// assert that the expectations were met
	mockGateway.AssertNotCalled(t, "AddUserPoolToAppSync", mock.Anything)
	mockGateway.AssertNotCalled(t, "GetUser", mock.Anything)
	mockGateway.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestAddGetUserFailure(t *testing.T) {
	testUserPoolID := aws.String("us-west-2_ZlG7Ldp1K")
	testAppClientID := aws.String("abcdefghijklmnopq")
	testIdentityPoolID := aws.String("us-west-2:abcdefghijklmnopq")
	testUserPool := &gateway.UserPool{
		UserPoolID:     testUserPoolID,
		AppClientID:    testAppClientID,
		IdentityPoolID: testIdentityPoolID,
	}
	testUserID := aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766")
	testGroup := aws.String("Admins")

	// create an instance of our test object
	mockGateway := &gateway.MockUserGateway{}
	// replace the global variable with our mock object
	userGateway = mockGateway

	// setup expectations
	mockGateway.On("CreateUserPool", testInput.DisplayName).Return(testUserPool, nil)
	mockGateway.On("CreateUserPoolGroups", testUserPoolID).Return(nil)
	mockGateway.On("CreateUser", testCreateUserInput).Return(testUserID, nil)
	mockGateway.On("AddUserToGroup", testUserID, testGroup, testUserPoolID).Return(nil)
	mockGateway.On("GetUser", testUserID, testUserPoolID).Return(&models.User{}, &genericapi.AWSError{})

	// call the code we are testing
	result, err := (API{}).CreateUserInfrastructure(testInput)

	// assert that the expectations were met
	mockGateway.AssertExpectations(t)
	assert.Nil(t, result)
	assert.Error(t, err)
}
