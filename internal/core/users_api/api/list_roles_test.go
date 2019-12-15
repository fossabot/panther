package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"

	"github.com/panther-labs/panther/api/lambda/users/models"
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
