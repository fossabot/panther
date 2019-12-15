package gateway

import (
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// DeleteUser calls cognito api delete user from a user pool
func (g *UsersGateway) DeleteUser(id *string, userPoolID *string) error {
	if _, err := g.userPoolClient.AdminDeleteUser(&provider.AdminDeleteUserInput{
		Username:   id,
		UserPoolId: userPoolID,
	}); err != nil {
		return &genericapi.AWSError{Method: "cognito.AdminDeleteUser", Err: err}
	}

	return nil
}
