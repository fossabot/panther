package api

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
