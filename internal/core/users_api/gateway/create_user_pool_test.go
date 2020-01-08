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
	fedProvider "github.com/aws/aws-sdk-go/service/cognitoidentity"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/pkg/genericapi"
)

var policyDoc = "%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22cognito-identity.amazonaws.com%22%7D%2C%22Action%22%3A%22sts%3AAssumeRoleWithWebIdentity%22%7D%5D%7D" //nolint:lll

func TestCreateUserPool(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	mockFedIdentityClient := &MockFedIdentityClient{}
	gw := &UsersGateway{
		userPoolClient:    mockCognitoClient,
		iamService:        mockIamService,
		fedIdentityClient: mockFedIdentityClient,
	}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn:                      aws.String("fakeArn"),
			AssumeRolePolicyDocument: aws.String(policyDoc),
		},
	}, nil)
	mockIamService.On(
		"UpdateAssumeRolePolicy", mock.Anything).Return(&iam.UpdateAssumeRolePolicyOutput{}, nil)
	mockCognitoClient.On("CreateUserPool", mock.Anything).Return(&provider.CreateUserPoolOutput{
		UserPool: &provider.UserPoolType{
			Id: aws.String("fakeuserpoolId"),
		},
	}, nil)
	mockCognitoClient.On("CreateUserPoolClient", mock.Anything).Return(
		&provider.CreateUserPoolClientOutput{
			UserPoolClient: &provider.UserPoolClientType{
				ClientId: aws.String("clientId"),
			},
		}, nil)
	mockFedIdentityClient.On("CreateIdentityPool", mock.Anything).Return(
		&fedProvider.IdentityPool{
			IdentityPoolId: aws.String("us-west-2:fakeidentitypoolid"),
			CognitoIdentityProviders: []*fedProvider.Provider{
				{
					ProviderName: aws.String("cognito-idp.us-west-2.amazonaws.com/fakeuserpoolId"),
				},
			},
		}, nil)

	mockFedIdentityClient.On("SetIdentityPoolRoles", mock.Anything).Return(
		&fedProvider.SetIdentityPoolRolesOutput{}, nil)

	mockCognitoClient.On("SetUserPoolMfaConfig", mock.Anything).Return(&provider.SetUserPoolMfaConfigOutput{}, nil)

	response, err := gw.CreateUserPool(aws.String("myAwesomeBirdCompany"))
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, &UserPool{
		AppClientID:    aws.String("clientId"),
		UserPoolID:     aws.String("fakeuserpoolId"),
		IdentityPoolID: aws.String("us-west-2:fakeidentitypoolid"),
	}, response)
	mockIamService.AssertExpectations(t)
	mockCognitoClient.AssertExpectations(t)
}

func TestCreatePoolGetRoleFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{}, &genericapi.AWSError{})

	response, err := gw.CreateUserPool(aws.String("myAwesomeBirdCompany"))
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
	mockIamService.AssertExpectations(t)
}

func TestCreatePoolFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)

	mockCognitoClient.On("CreateUserPool", mock.Anything).Return(&provider.CreateUserPoolOutput{}, &genericapi.AWSError{})

	response, err := gw.CreateUserPool(aws.String("myAwesomeBirdCompany"))
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
	mockIamService.AssertExpectations(t)
	mockCognitoClient.AssertExpectations(t)
}

func TestCreatePoolAppClientFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn:                      aws.String("fakeArn"),
			AssumeRolePolicyDocument: aws.String(policyDoc),
		},
	}, nil)

	mockIamService.On(
		"UpdateAssumeRolePolicy", mock.Anything).Return(&iam.UpdateAssumeRolePolicyOutput{}, nil)

	mockCognitoClient.On("CreateUserPool", mock.Anything).Return(&provider.CreateUserPoolOutput{
		UserPool: &provider.UserPoolType{
			Id: aws.String("fakeuserpoolId"),
		},
	}, nil)
	mockCognitoClient.On("CreateUserPoolClient", mock.Anything).Return(&provider.CreateUserPoolClientOutput{}, &genericapi.AWSError{})

	response, err := gw.CreateUserPool(aws.String("myAwesomeBirdCompany"))
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
	mockCognitoClient.AssertExpectations(t)
}

func TestCreateUserPoolIdentityPoolFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	mockFedIdentityClient := &MockFedIdentityClient{}
	gw := &UsersGateway{
		userPoolClient:    mockCognitoClient,
		iamService:        mockIamService,
		fedIdentityClient: mockFedIdentityClient,
	}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)

	mockCognitoClient.On("CreateUserPool", mock.Anything).Return(&provider.CreateUserPoolOutput{
		UserPool: &provider.UserPoolType{
			Id: aws.String("fakeuserpoolId"),
		},
	}, nil)
	mockCognitoClient.On("CreateUserPoolClient", mock.Anything).Return(
		&provider.CreateUserPoolClientOutput{
			UserPoolClient: &provider.UserPoolClientType{
				ClientId: aws.String("clientId"),
			},
		}, nil)
	mockFedIdentityClient.On("CreateIdentityPool", mock.Anything).Return(
		&fedProvider.IdentityPool{}, &genericapi.AWSError{})

	mockCognitoClient.On("SetUserPoolMfaConfig", mock.Anything).Return(&provider.SetUserPoolMfaConfigOutput{}, nil)

	response, err := gw.CreateUserPool(aws.String("myAwesomeBirdCompany"))
	assert.Nil(t, response)
	assert.Error(t, err)
}

func TestCreateUserPoolSetIdentityPoolFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	mockFedIdentityClient := &MockFedIdentityClient{}
	gw := &UsersGateway{
		userPoolClient:    mockCognitoClient,
		iamService:        mockIamService,
		fedIdentityClient: mockFedIdentityClient,
	}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)

	mockCognitoClient.On("CreateUserPool", mock.Anything).Return(&provider.CreateUserPoolOutput{
		UserPool: &provider.UserPoolType{
			Id: aws.String("fakeuserpoolId"),
		},
	}, nil)
	mockCognitoClient.On("CreateUserPoolClient", mock.Anything).Return(
		&provider.CreateUserPoolClientOutput{
			UserPoolClient: &provider.UserPoolClientType{
				ClientId: aws.String("clientId"),
			},
		}, nil)
	mockFedIdentityClient.On("CreateIdentityPool", mock.Anything).Return(
		&fedProvider.IdentityPool{
			IdentityPoolId: aws.String("us-west-2:fakeidentitypoolid"),
			CognitoIdentityProviders: []*fedProvider.Provider{
				{
					ProviderName: aws.String("cognito-idp.us-west-2.amazonaws.com/fakeuserpoolId"),
				},
			},
		}, nil)

	mockFedIdentityClient.On("SetIdentityPoolRoles", mock.Anything).Return(
		&fedProvider.SetIdentityPoolRolesOutput{}, &genericapi.AWSError{})

	mockCognitoClient.On("SetUserPoolMfaConfig", mock.Anything).Return(&provider.SetUserPoolMfaConfigOutput{}, nil)

	response, err := gw.CreateUserPool(aws.String("myAwesomeBirdCompany"))
	assert.Nil(t, response)
	assert.Error(t, err)
}

func TestCreatePoolMfaSettingsFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)

	mockCognitoClient.On("CreateUserPool", mock.Anything).Return(&provider.CreateUserPoolOutput{
		UserPool: &provider.UserPoolType{
			Id: aws.String("fakeuserpoolId"),
		},
	}, nil)
	mockCognitoClient.On("CreateUserPoolClient", mock.Anything).Return(&provider.CreateUserPoolClientOutput{}, nil)
	mockCognitoClient.On("SetUserPoolMfaConfig", mock.Anything).Return(&provider.SetUserPoolMfaConfigOutput{}, &genericapi.AWSError{})

	response, err := gw.CreateUserPool(aws.String("myAwesomeBirdCompany"))
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.IsType(t, err, &genericapi.AWSError{})
	mockIamService.AssertExpectations(t)
	mockCognitoClient.AssertExpectations(t)
}
