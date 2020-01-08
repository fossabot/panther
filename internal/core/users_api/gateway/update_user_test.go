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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/assert"
)

type mockUpdateUserClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockUpdateUserClient) AdminUpdateUserAttributes(
	*provider.AdminUpdateUserAttributesInput) (*provider.AdminUpdateUserAttributesOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}
	return &provider.AdminUpdateUserAttributesOutput{}, nil
}

func TestUpdateUserGivenName(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockUpdateUserClient{}}
	assert.NoError(t, gw.UpdateUser(&UpdateUserInput{
		GivenName:  aws.String("Richie"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("userPoolId"),
	}))
}

func TestUpdateUserFamilyName(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockUpdateUserClient{}}
	assert.NoError(t, gw.UpdateUser(&UpdateUserInput{
		FamilyName: aws.String("Homie"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("userPoolId"),
	}))
}

func TestUpdateUserPhoneNumber(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockUpdateUserClient{}}
	assert.NoError(t, gw.UpdateUser(&UpdateUserInput{
		PhoneNumber: aws.String("+1234567890"),
		ID:          aws.String("user123"),
		UserPoolID:  aws.String("userPoolId"),
	}))
}

func TestUpdateUserEmail(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockUpdateUserClient{}}
	assert.NoError(t, gw.UpdateUser(&UpdateUserInput{
		Email:      aws.String("rich@homie.com"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("userPoolId"),
	}))
}

func TestUpdateMultipleAttributes(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockUpdateUserClient{}}
	assert.NoError(t, gw.UpdateUser(&UpdateUserInput{
		GivenName:   aws.String("Richie"),
		FamilyName:  aws.String("Homie"),
		PhoneNumber: aws.String("+1234567890"),
		Email:       aws.String("rich@homie.com"),
		ID:          aws.String("user123"),
		UserPoolID:  aws.String("userPoolId"),
	}))
}

func TestUpdateUserFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockUpdateUserClient{serviceErr: true}}
	assert.Error(t, gw.UpdateUser(&UpdateUserInput{
		ID:         aws.String("user123"),
		GivenName:  aws.String("Richie"),
		UserPoolID: aws.String("userPoolId"),
	}))
}
