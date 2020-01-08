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
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	cfnIface "github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	fedProvider "github.com/aws/aws-sdk-go/service/cognitoidentity"
	fedProviderI "github.com/aws/aws-sdk-go/service/cognitoidentity/cognitoidentityiface"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/mock"
)

// ListGraphqlApis mocks AppSync.ListGraphqlApis
func (m *MockAppSyncService) ListGraphqlApis(input *appsync.ListGraphqlApisInput) (*appsync.ListGraphqlApisOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*appsync.ListGraphqlApisOutput), args.Error(1)
}

// MockCloudFormationService can be passed as a mock object to unit tests
type MockCloudFormationService struct {
	cfnIface.CloudFormationAPI
	mock.Mock
}

// DescribeStackResource mocks CFN.DescribeStackResource
func (m *MockCloudFormationService) DescribeStackResource(input *cloudformation.DescribeStackResourceInput) (
	*cloudformation.DescribeStackResourceOutput, error) {

	args := m.Called(input)
	return args.Get(0).(*cloudformation.DescribeStackResourceOutput), args.Error(1)
}

// MockIamService can be passed as a mock object to unit tests
type MockIamService struct {
	IAMService
	mock.Mock
}

// GetRole mocks IAM.GetRole
func (m *MockIamService) GetRole(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*iam.GetRoleOutput), args.Error(1)
}

// UpdateAssumeRolePolicy mocks IAM.UpdateAssumeRolePolicy
func (m *MockIamService) UpdateAssumeRolePolicy(input *iam.UpdateAssumeRolePolicyInput) (*iam.UpdateAssumeRolePolicyOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*iam.UpdateAssumeRolePolicyOutput), args.Error(1)
}

// MockAppSyncService can be passed as a mock object to unit tests
type MockAppSyncService struct {
	cfnIface.CloudFormationAPI
	mock.Mock
}

// MockCognitoClient can be passed as a mock object to unit tests
type MockCognitoClient struct {
	providerI.CognitoIdentityProviderAPI
	mock.Mock
}

// AdminAddUserToGroup mocks AdminAddUserToGroup for testing
func (m *MockCognitoClient) AdminAddUserToGroup(
	input *provider.AdminAddUserToGroupInput) (*provider.AdminAddUserToGroupOutput, error) {

	args := m.Called(input)
	return args.Get(0).(*provider.AdminAddUserToGroupOutput), args.Error(1)
}

// AdminCreateUser mocks AdminCreateUser for testing
func (m *MockCognitoClient) AdminCreateUser(
	input *provider.AdminCreateUserInput) (*provider.AdminCreateUserOutput, error) {

	args := m.Called(input)
	return args.Get(0).(*provider.AdminCreateUserOutput), args.Error(1)
}

// CreateGroup mocks CreateGroup for testing
func (m *MockCognitoClient) CreateGroup(
	input *provider.CreateGroupInput) (*provider.CreateGroupOutput, error) {

	args := m.Called(input)
	return args.Get(0).(*provider.CreateGroupOutput), args.Error(1)
}

// CreateUserPool mocks CreateUserPool for testing
func (m *MockCognitoClient) CreateUserPool(input *provider.CreateUserPoolInput) (*provider.CreateUserPoolOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*provider.CreateUserPoolOutput), args.Error(1)
}

// CreateUserPoolClient mocks CreateUserPoolClient for testing
func (m *MockCognitoClient) CreateUserPoolClient(input *provider.CreateUserPoolClientInput) (*provider.CreateUserPoolClientOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*provider.CreateUserPoolClientOutput), args.Error(1)
}

// SetUserPoolMfaConfig mocks SetUserPoolMfaConfig for testing
func (m *MockCognitoClient) SetUserPoolMfaConfig(input *provider.SetUserPoolMfaConfigInput) (*provider.SetUserPoolMfaConfigOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*provider.SetUserPoolMfaConfigOutput), args.Error(1)
}

// MockFedIdentityClient can be passed as a mock object to unit tests
type MockFedIdentityClient struct {
	fedProviderI.CognitoIdentityAPI
	mock.Mock
}

// GetCredentialsForIdentity mocks GetCredentialsForIdentity for testing
func (m *MockFedIdentityClient) GetCredentialsForIdentity(input *fedProvider.GetCredentialsForIdentityInput) (
	*fedProvider.GetCredentialsForIdentityOutput, error) {

	args := m.Called(input)
	return args.Get(0).(*fedProvider.GetCredentialsForIdentityOutput), args.Error(1)
}

// CreateIdentityPool mocks CreateIdentityPool for testing
func (m *MockFedIdentityClient) CreateIdentityPool(input *fedProvider.CreateIdentityPoolInput) (*fedProvider.IdentityPool, error) {
	args := m.Called(input)
	return args.Get(0).(*fedProvider.IdentityPool), args.Error(1)
}

// SetIdentityPoolRoles mocks SetIdentityPoolRoles for testing
func (m *MockFedIdentityClient) SetIdentityPoolRoles(input *fedProvider.SetIdentityPoolRolesInput) (
	*fedProvider.SetIdentityPoolRolesOutput, error) {

	args := m.Called(input)
	return args.Get(0).(*fedProvider.SetIdentityPoolRolesOutput), args.Error(1)
}
