package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
)

// CreateUserInfrastructure creates user infrastructure needed for a new organization
func (API) CreateUserInfrastructure(input *models.CreateUserInfrastructureInput) (*models.CreateUserInfrastructureOutput, error) {
	userPool, err := userGateway.CreateUserPool(input.DisplayName)
	if err != nil {
		zap.L().Error("error creating user pool", zap.Error(err))
		return nil, err
	}
	userPoolID := userPool.UserPoolID

	err = userGateway.CreateUserPoolGroups(userPoolID)
	if err != nil {
		zap.L().Error("error creating user pool groups", zap.Error(err))
		return nil, err
	}

	firstUserID, err := userGateway.CreateUser(&gateway.CreateUserInput{
		GivenName:  input.GivenName,
		FamilyName: input.FamilyName,
		Email:      input.Email,
		UserPoolID: userPoolID,
	})
	if err != nil {
		zap.L().Error("error creating first user", zap.Error(err))
		return nil, err
	}

	err = userGateway.AddUserToGroup(firstUserID, aws.String("Admins"), userPoolID)
	if err != nil {
		zap.L().Error("error adding user to Admin group", zap.Error(err))
		return nil, err
	}

	firstUser, err := userGateway.GetUser(firstUserID, userPoolID)
	if err != nil {
		zap.L().Error("error getting first user", zap.Error(err))
		return nil, err
	}

	return &models.CreateUserInfrastructureOutput{
		User:           firstUser,
		UserPoolID:     userPoolID,
		AppClientID:    userPool.AppClientID,
		IdentityPoolID: userPool.IdentityPoolID,
	}, nil
}
