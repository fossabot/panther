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

type mockGatewayUpdateUserClient struct {
	gateway.API
	updateErr bool
	listErr   bool
	removeErr bool
}

func (m *mockGatewayUpdateUserClient) UpdateUser(*gateway.UpdateUserInput) error {
	if m.updateErr {
		return &genericapi.AWSError{}
	}
	return nil
}

func (m *mockGatewayUpdateUserClient) ListGroupsForUser(*string, *string) ([]*models.Group, error) {
	if m.listErr {
		return nil, &genericapi.AWSError{}
	}
	return []*models.Group{
		{
			Name:        aws.String("Admins"),
			Description: aws.String("Administrator of the group"),
		},
	}, nil
}

func (m *mockGatewayUpdateUserClient) RemoveUserFromGroup(*string, *string, *string) error {
	if m.removeErr {
		return &genericapi.AWSError{}
	}
	return nil
}

func (m *mockGatewayUpdateUserClient) AddUserToGroup(*string, *string, *string) error {
	if m.removeErr {
		return &genericapi.AWSError{}
	}
	return nil
}

func TestUpdateUserGatewayErr(t *testing.T) {
	userGateway = &mockGatewayUpdateUserClient{updateErr: true}
	input := &models.UpdateUserInput{
		GivenName:  aws.String("Richie"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.Error(t, (API{}).UpdateUser(input))
}

func TestUpdateUserChangeRole(t *testing.T) {
	userGateway = &mockGatewayUpdateUserClient{}
	input := &models.UpdateUserInput{
		Role:       aws.String("Admins"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.NoError(t, (API{}).UpdateUser(input))
}

func TestUpdateUserChangeRoleListErr(t *testing.T) {
	userGateway = &mockGatewayUpdateUserClient{listErr: true}
	input := &models.UpdateUserInput{
		Role:       aws.String("Admins"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.Error(t, (API{}).UpdateUser(input))
}

func TestUpdateUserChangeRoleRemoveErr(t *testing.T) {
	userGateway = &mockGatewayUpdateUserClient{removeErr: true}
	input := &models.UpdateUserInput{
		Role:       aws.String("Admins"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.Error(t, (API{}).UpdateUser(input))
}

func TestUpdateUserHandle(t *testing.T) {
	userGateway = &mockGatewayUpdateUserClient{}
	input := &models.UpdateUserInput{
		GivenName:  aws.String("Richie"),
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.NoError(t, (API{}).UpdateUser(input))
}
