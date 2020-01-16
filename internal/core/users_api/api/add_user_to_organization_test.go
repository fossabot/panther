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
