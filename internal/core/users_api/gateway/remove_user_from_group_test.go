package gateway

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/assert"
)

type mockRemoveUserFromGroupClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockRemoveUserFromGroupClient) AdminRemoveUserFromGroup(
	*provider.AdminRemoveUserFromGroupInput) (*provider.AdminRemoveUserFromGroupOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}
	return &provider.AdminRemoveUserFromGroupOutput{}, nil
}

func TestRemoveUserFromGroup(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockRemoveUserFromGroupClient{}}
	assert.NoError(t, gw.RemoveUserFromGroup(
		aws.String("user123"),
		aws.String("Admins"),
		aws.String("userPoolId"),
	))
}

func TestRemoveUserFromGroupFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockRemoveUserFromGroupClient{serviceErr: true}}
	assert.Error(t, gw.RemoveUserFromGroup(
		aws.String("user123"),
		aws.String("Admins"),
		aws.String("userPoolId"),
	))
}
