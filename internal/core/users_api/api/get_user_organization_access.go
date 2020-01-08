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

	organizationmodels "github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// GetUserOrganizationAccess calls dynamodb to get user's organization id.
func (API) GetUserOrganizationAccess(input *models.GetUserOrganizationAccessInput) (*models.GetUserOrganizationAccessOutput, error) {
	// Delete user from Dynamo
	_, err := userTable.Get(input.Email)
	if err != nil {
		zap.L().Error("error getting user", zap.Error(err))
		return nil, err
	}
	org, err := GetOrganizations()
	if err != nil {
		zap.L().Error("error getting organization", zap.Error(err))
		return nil, err
	}
	return org, nil
}

// GetOrganizations calls the organization api to fetch access related identifiers
func GetOrganizations() (*models.GetUserOrganizationAccessOutput, error) {
	input := organizationmodels.LambdaInput{GetOrganization: &organizationmodels.GetOrganizationInput{}}
	var org organizationmodels.GetOrganizationOutput
	if err := genericapi.Invoke(lambdaClient, organizationAPI, &input, &org); err != nil {
		return nil, err
	}
	return &models.GetUserOrganizationAccessOutput{
		UserPoolID:     org.Organization.AwsConfig.UserPoolID,
		AppClientID:    org.Organization.AwsConfig.AppClientID,
		IdentityPoolID: org.Organization.AwsConfig.IdentityPoolID,
	}, nil
}
