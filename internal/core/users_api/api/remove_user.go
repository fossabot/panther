package api

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
