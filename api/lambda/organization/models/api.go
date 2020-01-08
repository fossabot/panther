package models

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

// LambdaInput is the request structure for the organization-api Lambda function.
type LambdaInput struct {
	CompleteAction     *CompleteActionInput     `json:"getCompletedActions"`
	CreateOrganization *CreateOrganizationInput `json:"createOrganization"`
	GetOrganization    *GetOrganizationInput    `json:"getOrganization"`
	UpdateOrganization *UpdateOrganizationInput `json:"updateOrganization"`
}

// CompleteActionInput Adds a Action to an Organization
type CompleteActionInput struct {
	CompletedActions []*Action `json:"actions"`
}

// CompleteActionOutput Adds a Action to an Organization
type CompleteActionOutput struct {
	CompletedActions []*Action `json:"actions"`
}

// CreateOrganizationInput creates a new Panther customer account.
type CreateOrganizationInput struct {
	AlertReportFrequency *string            `json:"alertReportFrequency" validate:"omitempty,oneof=P1D P1W"`
	AwsConfig            *AwsConfig         `json:"awsConfig"`
	DisplayName          *string            `json:"displayName" validate:"required,min=1"`
	Email                *string            `genericapi:"redact" json:"email" validate:"required,email"`
	Phone                *string            `genericapi:"redact" json:"phone"`
	RemediationConfig    *RemediationConfig `json:"remediationConfig,omitempty"`
}

// CreateOrganizationOutput returns the newly created organization.
type CreateOrganizationOutput struct {
	Organization *Organization `json:"organization"`
}

// GetOrganizationInput retrieves the details of a Panther customer account.
type GetOrganizationInput struct {
}

// GetOrganizationOutput is the table row representing a customer account.
type GetOrganizationOutput struct {
	Organization *Organization `json:"organization"`
}

// UpdateOrganizationInput modifies the details of an existing organization.
type UpdateOrganizationInput struct {
	CreateOrganizationInput
}

// UpdateOrganizationOutput is the table row representing the modified customer account.
type UpdateOrganizationOutput struct {
	Organization *Organization `json:"organization"`
}
