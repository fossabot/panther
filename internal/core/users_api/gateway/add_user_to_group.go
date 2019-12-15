package gateway

import (
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// AddUserToGroup calls cognito api add a user to a specified group
func (g *UsersGateway) AddUserToGroup(id *string, groupName *string, userPoolID *string) error {
	if _, err := g.userPoolClient.AdminAddUserToGroup(&provider.AdminAddUserToGroupInput{
		GroupName:  groupName,
		Username:   id,
		UserPoolId: userPoolID,
	}); err != nil {
		return &genericapi.AWSError{Method: "cognito.AdminAddUserToGroup", Err: err}
	}

	return nil
}
