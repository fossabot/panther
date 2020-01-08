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
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// RemoveUser deletes a user from cognito.
func (API) RemoveUser(input *models.RemoveUserInput) error {
	// Get user sub from Cognito
	user, err := userGateway.GetUser(input.ID, input.UserPoolID)
	if err != nil {
		zap.L().Error("error getting user from user pool", zap.Error(err))
		return err
	}

	// Delete user from Cognito user pool
	err = userGateway.DeleteUser(input.ID, input.UserPoolID)
	if err != nil {
		zap.L().Error("error deleting user from user pool", zap.Error(err))
		return err
	}

	// Delete user from Dynamo
	err = userTable.Delete(user.Email)
	if err != nil {
		zap.L().Error("error deleting user from dynamo", zap.Error(err))
		return err
	}

	return nil
}
