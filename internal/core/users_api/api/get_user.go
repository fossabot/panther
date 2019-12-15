package api

import (
	"github.com/panther-labs/panther/api/lambda/users/models"
)

// GetUser calls userGateway to get user information.
func (API) GetUser(input *models.GetUserInput) (*models.GetUserOutput, error) {
	user, err := userGateway.GetUser(input.ID, input.UserPoolID)
	if err != nil {
		return nil, err
	}

	groups, err := userGateway.ListGroupsForUser(input.ID, input.UserPoolID)
	if err != nil {
		return nil, err
	}

	if len(groups) > 0 {
		user.Role = groups[0].Name
	}
	return user, nil
}
