package api

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
