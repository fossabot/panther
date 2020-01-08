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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"
)

var testAddUserToGroupInput = &provider.AdminAddUserToGroupInput{
	GroupName:  aws.String("Admin"),
	Username:   aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766"),
	UserPoolId: aws.String("us-west-2_ZlG7Ldp1K"),
}

func TestAddUserToGroup(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On(
		"AdminAddUserToGroup", testAddUserToGroupInput).Return(&provider.AdminAddUserToGroupOutput{}, nil)

	assert.NoError(t, gw.AddUserToGroup(
		testAddUserToGroupInput.Username,
		testAddUserToGroupInput.GroupName,
		testAddUserToGroupInput.UserPoolId,
	))
	mockCognitoClient.AssertExpectations(t)
}

func TestAddUserToGroupFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On("AdminAddUserToGroup", testAddUserToGroupInput).Return(
		&provider.AdminAddUserToGroupOutput{}, &genericapi.AWSError{})

	assert.Error(t, gw.AddUserToGroup(
		testAddUserToGroupInput.Username,
		testAddUserToGroupInput.GroupName,
		testAddUserToGroupInput.UserPoolId,
	))
	mockCognitoClient.AssertExpectations(t)
}
