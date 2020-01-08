package api

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
