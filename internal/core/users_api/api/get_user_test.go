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
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	"github.com/panther-labs/panther/pkg/genericapi"
)

type mockGatewayGetUserClient struct {
	gateway.API
	getUserGatewayErr    bool
	listGroupsGatewayErr bool
}

func (m *mockGatewayGetUserClient) GetUser(id *string, userPoolID *string) (*models.User, error) {
	if m.getUserGatewayErr {
		return nil, &genericapi.AWSError{}
	}
	return &models.User{
		GivenName:   aws.String("Joe"),
		FamilyName:  aws.String("Blow"),
		ID:          id,
		Email:       aws.String("joe@blow.com"),
		PhoneNumber: aws.String("+1234567890"),
		CreatedAt:   aws.Int64(1545442826),
		Status:      aws.String("CONFIRMED"),
	}, nil
}

func (m *mockGatewayGetUserClient) ListGroupsForUser(*string, *string) ([]*models.Group, error) {
	if m.listGroupsGatewayErr {
		return nil, &genericapi.AWSError{}
	}
	return []*models.Group{
		{
			Description: aws.String("Roles Description"),
			Name:        aws.String("Admins"),
		},
	}, nil
}

func TestGetUserGatewayErr(t *testing.T) {
	userGateway = &mockGatewayGetUserClient{getUserGatewayErr: true}
	result, err := (API{}).GetUser(&models.GetUserInput{
		UserPoolID: aws.String("fakePoolId"),
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetUserListGroupsForUserGatewayErr(t *testing.T) {
	userGateway = &mockGatewayGetUserClient{listGroupsGatewayErr: true}
	result, err := (API{}).GetUser(&models.GetUserInput{
		UserPoolID: aws.String("fakePoolId"),
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetUserHandle(t *testing.T) {
	userGateway = &mockGatewayGetUserClient{}
	result, err := (API{}).GetUser(&models.GetUserInput{
		UserPoolID: aws.String("fakePoolId"),
	})
	assert.NotNil(t, result)
	assert.NoError(t, err)
}
