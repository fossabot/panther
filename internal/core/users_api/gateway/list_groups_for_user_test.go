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

	"github.com/panther-labs/panther/api/lambda/users/models"
)

type mockListGroupsForUserClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockListGroupsForUserClient) AdminListGroupsForUser(
	*provider.AdminListGroupsForUserInput) (*provider.AdminListGroupsForUserOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}

	return &provider.AdminListGroupsForUserOutput{
		Groups: []*provider.GroupType{
			{
				CreationDate:     &time.Time{},
				Description:      aws.String("Roles Description"),
				GroupName:        aws.String("Admins"),
				LastModifiedDate: &time.Time{},
				Precedence:       aws.Int64(0),
				RoleArn:          aws.String("arn::1234567"),
				UserPoolId:       aws.String("Pool 123"),
			},
		},
	}, nil
}

func TestListGroupsForUser(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListGroupsForUserClient{}}
	result, err := gw.ListGroupsForUser(aws.String("user123"), aws.String("fakePoolId"))
	groups := []*models.Group{
		{
			Description: aws.String("Roles Description"),
			Name:        aws.String("Admins"),
		},
	}
	assert.Equal(t, groups, result)
	assert.NoError(t, err)
}

func TestListGroupsForUserFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListGroupsForUserClient{serviceErr: true}}
	result, err := gw.ListGroupsForUser(aws.String("user123"), aws.String("fakePoolId"))
	assert.Nil(t, result)
	assert.Error(t, err)
}
