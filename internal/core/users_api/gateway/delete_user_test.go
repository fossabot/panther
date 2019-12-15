package gateway

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/assert"
)

type mockDeleteUserClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockDeleteUserClient) AdminDeleteUser(
	*provider.AdminDeleteUserInput) (*provider.AdminDeleteUserOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}
	return &provider.AdminDeleteUserOutput{}, nil
}

func TestDeleteUser(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockDeleteUserClient{}}
	assert.NoError(t, gw.DeleteUser(aws.String("user123"), aws.String("userPoolId")))
}

func TestDeleteUserFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockDeleteUserClient{serviceErr: true}}
	assert.Error(t, gw.DeleteUser(aws.String("user123"), aws.String("userPoolId")))
}
