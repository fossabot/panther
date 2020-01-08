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

// Action defines an action the organization took
type Action = string

const (
	// VisitedOnboardingFlow defines when an organization visited the onboarding flow
	VisitedOnboardingFlow Action = "VISITED_ONBOARDING_FLOW"
)

// Organization defines the fields in the table row.
type Organization struct {
	AlertReportFrequency *string            `json:"alertReportFrequency"`
	AwsConfig            *AwsConfig         `json:"awsConfig"`
	CompletedActions     []*Action          `dynamodbav:"completedActions,omitempty,stringset" json:"completedActions"`
	CreatedAt            *string            `json:"createdAt"`
	DisplayName          *string            `json:"displayName"`
	Email                *string            `json:"email"`
	Phone                *string            `json:"phone"`
	RemediationConfig    *RemediationConfig `json:"remediationConfig,omitempty"`
}

// AwsConfig defines metadata related to AWS infrastructure for the organization
type AwsConfig struct {
	UserPoolID     *string `json:"userPoolId"`
	AppClientID    *string `json:"appClientId"`
	IdentityPoolID *string `json:"identityPoolId"`
}

// RemediationConfig contains information related to Remediation actions
type RemediationConfig struct {
	// Each organization will have one Lambda that is able to perform remediation for their AWS infrastructure.
	// This field contains the ARN for that Lambda.
	AwsRemediationLambdaArn *string `json:"awsRemediationLambdaArn,omitempty"`
}
