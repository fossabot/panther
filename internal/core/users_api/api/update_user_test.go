package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"

	"github.com/panther-labs/panther/api/lambda/users/models"
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
