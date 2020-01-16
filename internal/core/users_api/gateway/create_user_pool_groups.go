package gateway

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
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"go.uber.org/zap"
)

// CreateUserPoolGroups creates the basic groups in the provided user pool
func (g *UsersGateway) CreateUserPoolGroups(userPoolID *string) error {
	adminsRoleArn, err := getRoleArn(g, aws.String(IdentityPoolAuthenticatedAdminsRole))
	if err != nil {
		zap.L().Error("error getting admin group arn", zap.Error(err))
		return err
	}

	adminGroupInput := &provider.CreateGroupInput{
		Description: aws.String("Administrators for Panther web application"),
		GroupName:   aws.String("Admins"),
		Precedence:  aws.Int64(0),
		UserPoolId:  userPoolID,
		RoleArn:     adminsRoleArn,
	}
	_, err = g.userPoolClient.CreateGroup(adminGroupInput)
	if err != nil {
		zap.L().Error("error creating admin group", zap.Error(err))
		return err
	}

	analystsRoleArn, err := getRoleArn(g, aws.String(IdentityPoolAuthenticatedAnalystsRole))
	if err != nil {
		zap.L().Error("error getting analyst group arn", zap.Error(err))
		return err
	}
	analystGroupInput := &provider.CreateGroupInput{
		Description: aws.String("Analysts for Panther web application, allows rule, alert and configuration managing"),
		GroupName:   aws.String("Analysts"),
		Precedence:  aws.Int64(1),
		UserPoolId:  userPoolID,
		RoleArn:     analystsRoleArn,
	}
	_, err = g.userPoolClient.CreateGroup(analystGroupInput)
	if err != nil {
		zap.L().Error("error creating analyst group", zap.Error(err))
		return err
	}

	readonlyRoleArn, err := getRoleArn(g, aws.String(IdentityPoolAuthenticatedReadOnlyRole))
	if err != nil {
		zap.L().Error("error getting readonlys group arn", zap.Error(err))
		return err
	}
	readonlyGroupInput := &provider.CreateGroupInput{
		Description: aws.String("ReadOnly Group for Panther web application, only have access to Read Only operations"),
		GroupName:   aws.String("ReadOnly"),
		Precedence:  aws.Int64(2),
		UserPoolId:  userPoolID,
		RoleArn:     readonlyRoleArn,
	}
	_, err = g.userPoolClient.CreateGroup(readonlyGroupInput)
	if err != nil {
		zap.L().Error("error creating readonly group", zap.Error(err))
		return err
	}

	return nil
}
