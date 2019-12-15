package api

import (
	"github.com/panther-labs/panther/api/lambda/users/models"
)

// ListRoles requests the groups from cognito.
func (API) ListRoles(input *models.ListRolesInput) (*models.ListRolesOutput, error) {
	groups, err := userGateway.ListGroups(input.UserPoolID)
	if err != nil {
		return nil, err
	}
	return &models.ListRolesOutput{Roles: groups}, nil
}
