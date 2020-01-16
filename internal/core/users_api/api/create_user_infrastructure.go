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
