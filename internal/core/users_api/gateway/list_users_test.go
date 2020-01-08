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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/assert"
)

type mockListUsersClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockListUsersClient) ListUsers(
	*provider.ListUsersInput) (*provider.ListUsersOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}

	return &provider.ListUsersOutput{
		Users: []*provider.UserType{
			{
				Attributes: []*provider.AttributeType{
					{
						Name:  aws.String("given_name"),
						Value: aws.String("Joe"),
					},
					{
						Name:  aws.String("family_name"),
						Value: aws.String("Blow"),
					},
					{
						Name:  aws.String("email"),
						Value: aws.String("joe@blow.com"),
					},
					{
						Name:  aws.String("phone_number"),
						Value: aws.String("+1234567890"),
					},
				},
				Enabled:              aws.Bool(true),
				UserCreateDate:       &time.Time{},
				UserLastModifiedDate: &time.Time{},
				Username:             aws.String("user123"),
				UserStatus:           aws.String("CONFIRMED"),
			},
		},
		PaginationToken: aws.String("token123"),
	}, nil
}

func TestListUsers(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListUsersClient{}}
	result, err := gw.ListUsers(
		aws.Int64(10),
		aws.String("paginationToken"),
		aws.String("userPoolId"),
	)
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestListUsersFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListUsersClient{serviceErr: true}}
	result, err := gw.ListUsers(
		aws.Int64(10),
		aws.String("paginationToken"),
		aws.String("userPoolId"),
	)
	assert.Nil(t, result)
	assert.Error(t, err)
}
