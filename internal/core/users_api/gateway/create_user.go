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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// CreateUserInput is input for CreateUser request
type CreateUserInput struct {
	GivenName   *string `json:"givenName"`
	FamilyName  *string `json:"familyName"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	UserPoolID  *string `json:"userPoolId"`
}

// Create a AdminCreateUserInput from the CreateUserInput.
func (g *UsersGateway) cognitoInputFromAPIInput(
	input *CreateUserInput) *provider.AdminCreateUserInput {

	var userAttrs []*provider.AttributeType
	var lowercaseEmail = aws.String(strings.ToLower(*input.Email))
	if input.GivenName != nil {
		userAttrs = append(userAttrs, &provider.AttributeType{
			Name:  aws.String("given_name"),
			Value: input.GivenName,
		})
	}

	if input.FamilyName != nil {
		userAttrs = append(userAttrs, &provider.AttributeType{
			Name:  aws.String("family_name"),
			Value: input.FamilyName,
		})
	}

	if input.Email != nil {
		userAttrs = append(userAttrs, &provider.AttributeType{
			Name:  aws.String("email"),
			Value: lowercaseEmail,
		})
	}

	if input.PhoneNumber != nil {
		userAttrs = append(userAttrs, &provider.AttributeType{
			Name:  aws.String("phone_number"),
			Value: input.PhoneNumber,
		})
	}

	userAttrs = append(userAttrs, &provider.AttributeType{
		Name:  aws.String("email_verified"),
		Value: aws.String("true"),
	})

	return &provider.AdminCreateUserInput{
		DesiredDeliveryMediums: []*string{aws.String("EMAIL")}, // todo: Get from os environment or configuration database
		ForceAliasCreation:     aws.Bool(false),                // todo: Get from os environment or configuration database
		UserAttributes:         userAttrs,
		Username:               lowercaseEmail,
		UserPoolId:             input.UserPoolID,
	}
}

// CreateUser calls cognito api and creates a new user with specified attributes and sends out an email invitation
func (g *UsersGateway) CreateUser(input *CreateUserInput) (*string, error) {
	cognitoInput := g.cognitoInputFromAPIInput(input)
	userOutput, err := g.userPoolClient.AdminCreateUser(cognitoInput)
	if err != nil {
		return nil, &genericapi.AWSError{Method: "cognito.AdminCreateUser", Err: err}
	}
	return userOutput.User.Username, nil
}
