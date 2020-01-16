package api

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
