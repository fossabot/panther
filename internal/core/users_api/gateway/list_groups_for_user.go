package gateway

import (
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// ListGroupsForUser calls cognito api to list groups that a user belongs to
func (g *UsersGateway) ListGroupsForUser(id *string, userPoolID *string) ([]*models.Group, error) {
	o, err := g.userPoolClient.AdminListGroupsForUser(&provider.AdminListGroupsForUserInput{
		Username:   id,
		UserPoolId: userPoolID,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.AdminListGroupsForUser", Err: err}
	}

	groups := make([]*models.Group, len(o.Groups))
	for i, og := range o.Groups {
		groups[i] = &models.Group{Description: og.Description, Name: og.GroupName}
	}
	return groups, nil
}
