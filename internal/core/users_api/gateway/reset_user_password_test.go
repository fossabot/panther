package gateway

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	providerI "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/stretchr/testify/assert"
)

type mockResetUserPasswordClient struct {
	providerI.CognitoIdentityProviderAPI
	serviceErr bool
}

func (m *mockResetUserPasswordClient) AdminResetUserPassword(
	*provider.AdminResetUserPasswordInput) (*provider.AdminResetUserPasswordOutput, error) {

	if m.serviceErr {
		return nil, errors.New("cognito does not exist")
	}
	return &provider.AdminResetUserPasswordOutput{}, nil
}

func TestResetUserPassword(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockResetUserPasswordClient{}}
	assert.NoError(t, gw.ResetUserPassword(aws.String("user123"), aws.String("userPoolId")))
}

func TestResetUserPasswordFailed(t *testing.T) {
	gw := &UsersGateway{userPoolClient: &mockResetUserPasswordClient{serviceErr: true}}
	assert.Error(t, gw.ResetUserPassword(aws.String("user123"), aws.String("userPoolId")))
}
