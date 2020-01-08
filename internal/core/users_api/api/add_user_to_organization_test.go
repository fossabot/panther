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

	"github.com/panther-labs/panther/api/lambda/users/models"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
	"github.com/panther-labs/panther/pkg/genericapi"
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
