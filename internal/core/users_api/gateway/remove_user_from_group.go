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
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// RemoveUserFromGroup calls cognito api to remove a user from a specified group
func (g *UsersGateway) RemoveUserFromGroup(id *string, groupName *string, userPoolID *string) error {
	if _, err := g.userPoolClient.AdminRemoveUserFromGroup(&provider.AdminRemoveUserFromGroupInput{
		GroupName:  groupName,
		Username:   id,
		UserPoolId: userPoolID,
	}); err != nil {
		return &genericapi.AWSError{Method: "cognito.AdminRemoveUserFromGroup", Err: err}
	}

	return nil
}
