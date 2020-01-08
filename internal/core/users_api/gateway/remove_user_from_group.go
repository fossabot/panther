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
