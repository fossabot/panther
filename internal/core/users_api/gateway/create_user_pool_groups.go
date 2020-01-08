package gateway

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
