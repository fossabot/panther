package api

import (
	"github.com/panther-labs/panther/api/lambda/users/models"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// AddUserToOrganization adds a mapping of user to organization ID
func (API) AddUserToOrganization(
	input *models.AddUserToOrganizationInput) (*models.AddUserToOrganizationOutput, error) {

	err := AddUserToOrganization(userTable, &models.UserItem{
		ID: input.Email,
	})
	if err != nil {
		return nil, err
	}

	return &models.AddUserToOrganizationOutput{Email: input.Email}, nil
}

// AddUserToOrganization adds a user to organization mapping in the table if it does not already exist
func AddUserToOrganization(table users.API, item *models.UserItem) error {
	// Check if user is already mapped to an organization
	existingUser, err := table.Get(item.ID)
	if existingUser != nil {
		return &genericapi.AlreadyExistsError{Message: "user already exists: " + *item.ID}
	}

	// If it a does not exist error, that is expected so continue
	// If it is another error, then return the error
	if _, isMissing := err.(*genericapi.DoesNotExistError); !isMissing {
		return err
	}

	// Add user - org mapping
	return table.Put(item)
}
