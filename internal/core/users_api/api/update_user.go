package api

import (
	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
)

func changeUserRole(input *models.UpdateUserInput) error {
	groups, err := userGateway.ListGroupsForUser(input.ID, input.UserPoolID)
	if err != nil {
		return err
	}

	for _, group := range groups {
		if err = userGateway.RemoveUserFromGroup(input.ID, group.Name, input.UserPoolID); err != nil {
			return err
		}
	}

	return userGateway.AddUserToGroup(input.ID, input.Role, input.UserPoolID)
}

// UpdateUser modifies user attributes and roles.
func (API) UpdateUser(input *models.UpdateUserInput) error {
	// Update basic user attributes if needed.
	if input.GivenName != nil || input.FamilyName != nil || input.PhoneNumber != nil {
		if err := userGateway.UpdateUser(&gateway.UpdateUserInput{
			GivenName:   input.GivenName,
			FamilyName:  input.FamilyName,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
			ID:          input.ID,
			UserPoolID:  input.UserPoolID,
		}); err != nil {
			return err
		}
	}

	// Change role if needed.
	if input.Role != nil {
		return changeUserRole(input)
	}

	return nil
}
