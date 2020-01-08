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
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/api/lambda/organization/models"
)

// CreateOrganization generates a new organization ID.
//
// TODO - populate the rules table for new customers
func (API) CreateOrganization(
	input *models.CreateOrganizationInput) (*models.CreateOrganizationOutput, error) {

	// Then write the new org to the Dynamo table
	org := &models.Organization{
		AlertReportFrequency: input.AlertReportFrequency,
		AwsConfig:            input.AwsConfig,
		CreatedAt:            aws.String(time.Now().Format(time.RFC3339)),
		DisplayName:          input.DisplayName,
		Email:                input.Email,
		Phone:                input.Phone,
		RemediationConfig:    input.RemediationConfig,
	}

	if err := orgTable.Put(org); err != nil {
		return nil, err
	}
	return &models.CreateOrganizationOutput{Organization: org}, nil
}
