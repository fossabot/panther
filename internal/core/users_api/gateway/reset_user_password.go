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

// ResetUserPassword calls cognito api to reset user password
func (g *UsersGateway) ResetUserPassword(id *string, userPoolID *string) error {
	if _, err := g.userPoolClient.AdminResetUserPassword(&provider.AdminResetUserPasswordInput{
		Username:   id,
		UserPoolId: userPoolID,
	}); err != nil {
		return &genericapi.AWSError{Method: "cognito.AdminResetUserPassword", Err: err}
	}

	return nil
}
