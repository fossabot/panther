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
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	"github.com/panther-labs/panther/pkg/genericapi"
)

type mockGatewayListRolesClient struct {
	gateway.API
	gatewayErr bool
}

func (m *mockGatewayListRolesClient) ListGroups(*string) ([]*models.Group, error) {
	if m.gatewayErr {
		return nil, &genericapi.AWSError{}
	}

	return []*models.Group{
		{
			Name:        aws.String("Admins"),
			Description: aws.String("High and mighty ones"),
		},
	}, nil
}

func TestListRolesGatewayErr(t *testing.T) {
	userGateway = &mockGatewayListRolesClient{gatewayErr: true}
	result, err := (API{}).ListRoles(&models.ListRolesInput{
		UserPoolID: aws.String("fakePoolId"),
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestListRolesHandle(t *testing.T) {
	userGateway = &mockGatewayListRolesClient{}
	result, err := (API{}).ListRoles(&models.ListRolesInput{
		UserPoolID: aws.String("fakePoolId"),
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(result.Roles))
}
