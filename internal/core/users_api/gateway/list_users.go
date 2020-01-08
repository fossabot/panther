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

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ListUsersOutput is output type for ListUsers
type ListUsersOutput struct {
	Users           []*models.User
	PaginationToken *string
}

func mapCognitoUserTypeToUser(u *provider.UserType) *models.User {
	user := models.User{
		CreatedAt: aws.Int64(u.UserCreateDate.Unix()),
		ID:        u.Username,
		Status:    u.UserStatus,
	}

	for _, attribute := range u.Attributes {
		switch *attribute.Name {
		case "email":
			user.Email = attribute.Value
		case "phone_number":
			user.PhoneNumber = attribute.Value
		case "given_name":
			user.GivenName = attribute.Value
		case "family_name":
			user.FamilyName = attribute.Value
		}
	}

	return &user
}

// ListUsers calls cognito api to list users that belongs to a user pool
func (g *UsersGateway) ListUsers(limit *int64, paginationToken *string, userPoolID *string) (*ListUsersOutput, error) {
	usersOutput, err := g.userPoolClient.ListUsers(&provider.ListUsersInput{
		Limit:           limit,
		PaginationToken: paginationToken,
		UserPoolId:      userPoolID,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.ListUsers", Err: err}
	}

	users := make([]*models.User, len(usersOutput.Users))
	for i, uo := range usersOutput.Users {
		users[i] = mapCognitoUserTypeToUser(uo)
	}
	return &ListUsersOutput{Users: users, PaginationToken: usersOutput.PaginationToken}, nil
}
