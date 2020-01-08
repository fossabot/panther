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

func mapGetUserOutputToPantherUser(u *provider.AdminGetUserOutput) *models.User {
	user := models.User{
		CreatedAt: aws.Int64(u.UserCreateDate.Unix()),
		ID:        u.Username,
		Status:    u.UserStatus,
	}

	for _, attribute := range u.UserAttributes {
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

// GetUser calls cognito api to get user info
func (g *UsersGateway) GetUser(id *string, userPoolID *string) (*models.User, error) {
	user, err := g.userPoolClient.AdminGetUser(&provider.AdminGetUserInput{
		Username:   id,
		UserPoolId: userPoolID,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.AdminGetUser", Err: err}
	}
	return mapGetUserOutputToPantherUser(user), nil
}
