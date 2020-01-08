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

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ListGroups calls cognito api to list groups for the user pool
func (g *UsersGateway) ListGroups(userPoolID *string) ([]*models.Group, error) {
	o, err := g.userPoolClient.ListGroups(&provider.ListGroupsInput{UserPoolId: userPoolID})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.ListGroups", Err: err}
	}

	groups := make([]*models.Group, len(o.Groups))
	for i, og := range o.Groups {
		groups[i] = &models.Group{Description: og.Description, Name: og.GroupName}
	}
	return groups, nil
}
