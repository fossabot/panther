package gateway

import (
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// ListGroups calls cognito api to list groups for the user pool
func (g *UsersGateway) ListGroups(userPoolID *string) ([]*models.Group, error) {
	o, err := g.userPoolClient.ListGroups(&provider.ListGroupsInput{UserPoolId: userPoolID})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.ListGroups", Err: err}
	}

	groups := make([]*models.Group, len(o.Groups))
	for i, og := range o.Groups {
		groups[i] = &models.Group{Description: og.Description, Name: og.GroupName}
	}
	return groups, nil
}
