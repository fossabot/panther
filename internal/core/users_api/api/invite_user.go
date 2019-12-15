package api

import (
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// InviteUser adds a new user to the Cognito user pool.
func (API) InviteUser(input *models.InviteUserInput) (*models.InviteUserOutput, error) {
	// Add user to org mapping in dynamo
	err := AddUserToOrganization(userTable, &models.UserItem{
		ID: input.Email,
	})
	if err != nil {
		return nil, err
	}

	// Create user in Cognito
	id, err := userGateway.CreateUser(&gateway.CreateUserInput{
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
		Email:      input.Email,
		UserPoolID: input.UserPoolID,
	})
	if err != nil {
		if deleteErr := userTable.Delete(input.Email); deleteErr != nil {
			zap.L().Error("error deleting user from dynamo after failed invitation", zap.Error(deleteErr))
		}
		return nil, err
	}

	if err = userGateway.AddUserToGroup(id, input.Role, input.UserPoolID); err != nil {
		return nil, err
	}

	return &models.InviteUserOutput{ID: id}, nil
}
