package gateway

import (
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// ResetUserPassword calls cognito api to reset user password
func (g *UsersGateway) ResetUserPassword(id *string, userPoolID *string) error {
	if _, err := g.userPoolClient.AdminResetUserPassword(&provider.AdminResetUserPasswordInput{
		Username:   id,
		UserPoolId: userPoolID,
	}); err != nil {
		return &genericapi.AWSError{Method: "cognito.AdminResetUserPassword", Err: err}
	}

	return nil
}
