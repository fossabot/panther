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
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/pkg/genericapi"
)

func TestCreateGroups(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)
	mockCognitoClient.On(
		"CreateGroup", mock.Anything).Return(&provider.CreateGroupOutput{}, nil)

	assert.NoError(t, gw.CreateUserPoolGroups(aws.String("us-west-2_ZlG7Ldp1K")))
	mockCognitoClient.AssertExpectations(t)
	mockCognitoClient.AssertNumberOfCalls(t, "CreateGroup", 3)
}

func TestCreateGroupsFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)
	mockCognitoClient.On(
		"CreateGroup", mock.Anything).Return(&provider.CreateGroupOutput{}, &genericapi.AWSError{})

	assert.Error(t, gw.CreateUserPoolGroups(aws.String("us-west-2_ZlG7Ldp1K")))
	mockCognitoClient.AssertExpectations(t)
	mockCognitoClient.AssertNumberOfCalls(t, "CreateGroup", 1)
}

func TestGetAdminRoleFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}
	mockRoleParam := iam.GetRoleInput{
		RoleName: aws.String(IdentityPoolAuthenticatedAdminsRole),
	}
	mockIamService.On(
		"GetRole", &mockRoleParam).Return(&iam.GetRoleOutput{}, &genericapi.AWSError{})

	assert.Error(t, gw.CreateUserPoolGroups(aws.String("us-west-2_ZlG7Ldp1K")))
	mockCognitoClient.AssertExpectations(t)
	mockCognitoClient.AssertNumberOfCalls(t, "CreateGroup", 0)
}
