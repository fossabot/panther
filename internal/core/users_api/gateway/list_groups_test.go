package gateway

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

type mockListGroupsClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockListGroupsClient) ListGroups(
	*provider.ListGroupsInput) (*provider.ListGroupsOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}

	groups := []*provider.GroupType{
		{
			CreationDate:     &time.Time{},
			Description:      aws.String("Roles Description"),
			GroupName:        aws.String("Admins"),
			LastModifiedDate: &time.Time{},
			Precedence:       aws.Int64(0),
			RoleArn:          aws.String("arn::1234567"),
			UserPoolId:       aws.String("Pool 123"),
		},
	}
	return &provider.ListGroupsOutput{Groups: groups}, nil
}

func TestListGroups(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListGroupsClient{}}
	result, err := gw.ListGroups(
		aws.String("fakePoolId"),
	)
	groups := []*models.Group{
		{
			Description: aws.String("Roles Description"),
			Name:        aws.String("Admins"),
		},
	}
	assert.Equal(t, groups, result)
	assert.NoError(t, err)
}

func TestListGroupsFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockListGroupsClient{serviceErr: true}}
	result, err := gw.ListGroups(
		aws.String("fakePoolId"),
	)
	assert.Nil(t, result)
	assert.Error(t, err)
}
