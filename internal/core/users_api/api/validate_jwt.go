package api

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
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ValidateCredentials validate the JWT token and returns the claims
func (API) ValidateCredentials(input *models.ValidateCredentialsInput) (*models.ValidateCredentialsOutput, error) {
	jwtc, err := userGateway.ValidateToken(input.IdentityID, input.JWT)
	if err != nil {
		zap.L().Error("error parsing JWT token", zap.Error(err))
		return nil, err
	}
	// Check if the credentials from the token matches with what we got in the database
	o, err := GetOrganizations()
	// if of these checks below error, we ought to get paged because someone is messing with us
	if o == nil {
		zap.L().Error("organization not found", zap.Error(err))
		return nil, &genericapi.InvalidInputError{Message: "organization not found"}
	}
	if *o.IdentityPoolID != *input.IdentityPoolID {
		zap.L().Error("identity Pool id to validate: " +
			*input.IdentityPoolID + " and identity pool id stored: " +
			*o.IdentityPoolID + " not matching")
		return nil, &genericapi.InvalidInputError{Message: "identity Pool Invalid"}
	}
	return &models.ValidateCredentialsOutput{
		Identity: jwtc,
	}, nil
}
