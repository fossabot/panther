package api

import (
	"github.com/panther-labs/panther/api/lambda/users/models"
)

// ListUsers lists details for each user in Panther.
func (API) ListUsers(input *models.ListUsersInput) (*models.ListUsersOutput, error) {
	listOutput, err := userGateway.ListUsers(input.Limit, input.PaginationToken, input.UserPoolID)
	if err != nil {
		return nil, err
	}

	for _, user := range listOutput.Users {
		groups, err := userGateway.ListGroupsForUser(user.ID, input.UserPoolID)
		if err != nil {
			return nil, err
		}
		if len(groups) > 0 {
			user.Role = groups[0].Name
		}
	}

	return &models.ListUsersOutput{
		Users:           listOutput.Users,
		PaginationToken: listOutput.PaginationToken,
	}, nil
}
