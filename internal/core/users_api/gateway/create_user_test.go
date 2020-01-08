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

var testCreateUserInput = &CreateUserInput{
	GivenName:   aws.String("Joe"),
	FamilyName:  aws.String("Blow"),
	Email:       aws.String("joe.blow@toe.com"),
	PhoneNumber: aws.String("+11234567890"),
	UserPoolID:  aws.String("userPoolId"),
}

var testAdminCreateUserInput = &provider.AdminCreateUserInput{
	DesiredDeliveryMediums: []*string{aws.String("EMAIL")},
	ForceAliasCreation:     aws.Bool(false),
	UserAttributes: []*provider.AttributeType{
		{
			Name:  aws.String("given_name"),
			Value: testCreateUserInput.GivenName,
		},
		{
			Name:  aws.String("family_name"),
			Value: testCreateUserInput.FamilyName,
		},
		{
			Name:  aws.String("email"),
			Value: testCreateUserInput.Email,
		},
		{
			Name:  aws.String("phone_number"),
			Value: testCreateUserInput.PhoneNumber,
		},
		{
			Name:  aws.String("email_verified"),
			Value: aws.String("true"),
		},
	},
	Username:   testCreateUserInput.Email,
	UserPoolId: testCreateUserInput.UserPoolID,
}

func TestCreateUser(t *testing.T) {
	testUserID := aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766")
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On(
		"AdminCreateUser", testAdminCreateUserInput).Return(&provider.AdminCreateUserOutput{
		User: &provider.UserType{
			Username: testUserID,
		},
	}, nil)

	id, err := gw.CreateUser(testCreateUserInput)

	assert.Equal(t, id, testUserID)
	assert.NoError(t, err)
	mockCognitoClient.AssertExpectations(t)
}

func TestCreateUserFailed(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On("AdminCreateUser", testAdminCreateUserInput).Return(
		&provider.AdminCreateUserOutput{}, &genericapi.AWSError{})

	id, err := gw.CreateUser(testCreateUserInput)

	assert.Nil(t, id)
	assert.Error(t, err)
	mockCognitoClient.AssertExpectations(t)
}
