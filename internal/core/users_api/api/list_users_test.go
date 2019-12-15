package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

type mockGatewayListUsersClient struct {
	gateway.API
	listUserGatewayErr  bool
	listGroupGatewayErr bool
}

func (m *mockGatewayListUsersClient) ListUsers(
	limit *int64, paginationToken *string, userPoolID *string) (*gateway.ListUsersOutput, error) {

	if m.listUserGatewayErr {
		return nil, &genericapi.AWSError{}
	}

	return &gateway.ListUsersOutput{
		Users: []*models.User{
			{
				GivenName:   aws.String("Joe"),
				FamilyName:  aws.String("Blow"),
				ID:          aws.String("user123"),
				Email:       aws.String("joe@blow.com"),
				PhoneNumber: aws.String("+1234567890"),
				CreatedAt:   aws.Int64(1545442826),
				Status:      aws.String("CONFIRMED"),
			},
		},
		PaginationToken: paginationToken,
	}, nil
}

func (m *mockGatewayListUsersClient) ListGroupsForUser(*string, *string) ([]*models.Group, error) {
	if m.listGroupGatewayErr {
		return nil, &genericapi.AWSError{}
	}

	return []*models.Group{
		{
			Description: aws.String("Roles Description"),
			Name:        aws.String("Admins"),
		},
	}, nil
}

func TestListUsersGatewayErr(t *testing.T) {
	userGateway = &mockGatewayListUsersClient{listUserGatewayErr: true}
	result, err := (API{}).ListUsers(&models.ListUsersInput{
		UserPoolID:      aws.String("fakePoolId"),
		Limit:           aws.Int64(10),
		PaginationToken: aws.String("paginationToken"),
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestListGroupsForUserGatewayErr(t *testing.T) {
	userGateway = &mockGatewayListUsersClient{listGroupGatewayErr: true}
	result, err := (API{}).ListUsers(&models.ListUsersInput{
		UserPoolID:      aws.String("fakePoolId"),
		Limit:           aws.Int64(10),
		PaginationToken: aws.String("paginationToken"),
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestListUsersHandle(t *testing.T) {
	userGateway = &mockGatewayListUsersClient{}
	result, err := (API{}).ListUsers(&models.ListUsersInput{
		UserPoolID:      aws.String("fakePoolId"),
		Limit:           aws.Int64(10),
		PaginationToken: aws.String("paginationToken"),
	})
	assert.NotNil(t, result)
	assert.NoError(t, err)
}
