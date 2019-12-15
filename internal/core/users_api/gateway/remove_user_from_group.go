package gateway

import (
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// RemoveUserFromGroup calls cognito api to remove a user from a specified group
func (g *UsersGateway) RemoveUserFromGroup(id *string, groupName *string, userPoolID *string) error {
	if _, err := g.userPoolClient.AdminRemoveUserFromGroup(&provider.AdminRemoveUserFromGroupInput{
		GroupName:  groupName,
		Username:   id,
		UserPoolId: userPoolID,
	}); err != nil {
		return &genericapi.AWSError{Method: "cognito.AdminRemoveUserFromGroup", Err: err}
	}

	return nil
}
