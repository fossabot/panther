package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	"github.com/panther-labs/panther/pkg/genericapi"
)

type mockGatewayResetUserPasswordClient struct {
	gateway.API
	gatewayErr bool
}

func (m *mockGatewayResetUserPasswordClient) ResetUserPassword(*string, *string) error {
	if m.gatewayErr {
		return &genericapi.AWSError{}
	}
	return nil
}

func TestResetUserPasswordGatewayErr(t *testing.T) {
	userGateway = &mockGatewayResetUserPasswordClient{gatewayErr: true}
	input := &models.ResetUserPasswordInput{
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.Error(t, (API{}).ResetUserPassword(input))
}

func TestResetUserPasswordHandle(t *testing.T) {
	userGateway = &mockGatewayResetUserPasswordClient{}
	input := &models.ResetUserPasswordInput{
		ID:         aws.String("user123"),
		UserPoolID: aws.String("fakePoolId"),
	}
	assert.NoError(t, (API{}).ResetUserPassword(input))
}
