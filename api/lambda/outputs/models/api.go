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

// LambdaInput is the invocation event expected by the Lambda function.
//
// Exactly one action must be specified.
type LambdaInput struct {
	AddOutput              *AddOutputInput              `json:"addOutput"`
	UpdateOutput           *UpdateOutputInput           `json:"updateOutput"`
	GetOutput              *GetOutputInput              `json:"getOutput"`
	DeleteOutput           *DeleteOutputInput           `json:"deleteOutput"`
	GetOrganizationOutputs *GetOrganizationOutputsInput `json:"getOrganizationOutputs"`
	SetDefaultOutputs      *SetDefaultOutputsInput      `json:"setDefaultOutputs"`
	GetDefaultOutputs      *GetDefaultOutputsInput      `json:"getDefaultOutputs"`
}

// AddOutputInput adds a new encrypted alert output to DynamoDB.
//
// Example:
// {
//     "addOutput": {
//         "displayName": "alert-channel",
//         "userId": "f6cfad0a-9bb0-4681-9503-02c54cc979c7",
//         "slack": {
//             "webhookURL": "https://hooks.slack.com/services/..."
//         }
//     }
// }
type AddOutputInput struct {
	UserID             *string       `json:"userId" validate:"required,uuid4"`
	DisplayName        *string       `json:"displayName" validate:"required,min=1"`
	OutputConfig       *OutputConfig `json:"outputConfig" validate:"required"`
	DefaultForSeverity []*string     `json:"defaultForSeverity"`
}

// AddOutputOutput returns a randomly generated UUID for the output.
//
// Example:
// {
//     "displayName": "alert-channel",
//     "outputId": "7d1c5854-f3ea-491c-8a52-0aa0d58cb456",
//     "outputType": "slack"
// }
type AddOutputOutput = AlertOutput

// DeleteOutputInput permanently deletes output credentials.
//
// Example:
// {
//     "deleteOutput": {
//         "outputId": "7d1c5854-f3ea-491c-8a52-0aa0d58cb456"
//     }
// }
type DeleteOutputInput struct {
	OutputID *string `json:"outputId" validate:"required,uuid4"`
	Force    *bool   `json:"force"`
}

// UpdateOutputInput updates an alert output configuration.
//
// Example:
// {
//     "updateOutput": {
//         "userId": "9d1c5854-f3ea-491c-8a52-0aa0d58cb456",
//         "outputId": "7d1c5854-f3ea-491c-8a52-0aa0d58cb456"
//     }
// }
type UpdateOutputInput struct {
	UserID             *string       `json:"userId" validate:"required,uuid4"`
	DisplayName        *string       `json:"displayName" validate:"required,min=1"`
	OutputID           *string       `json:"outputId" validate:"required,uuid4"`
	OutputConfig       *OutputConfig `json:"outputConfig" validate:"required"`
	DefaultForSeverity []*string     `json:"defaultForSeverity"`
}

// UpdateOutputOutput returns the new updated output
//
// Example:
// {
//     "displayName": "alert-channel",
//     "outputId": "7d1c5854-f3ea-491c-8a52-0aa0d58cb456",
//     "outputType": "slack"
// }
type UpdateOutputOutput = AlertOutput

// GetOutputInput fetches the configuration for a specific alert output id of an organization
type GetOutputInput struct {
	OutputID *string `json:"outputId" validate:"required,uuid4"`
}

// GetOutputOutput contains the configuration for an alert
type GetOutputOutput = AlertOutput

// GetOrganizationOutputsInput fetches all alert output configuration for one organization
//
// Example:
// {
//     "getOrganizationOutputsInput": {
//     }
// }
type GetOrganizationOutputsInput struct {
}

// GetOrganizationOutputsOutput returns all the alert outputs for one organization
//
// Example:
// {
//     "displayName": "alert-channel",
//     "outputId": "7d1c5854-f3ea-491c-8a52-0aa0d58cb456",
//     "outputType": "slack"
// }
type GetOrganizationOutputsOutput = []*AlertOutput

// SetDefaultOutputsInput sets the default output for an organization
type SetDefaultOutputsInput struct {
	Severity  *string   `json:"severity" validate:"required,oneof=INFO LOW MEDIUM HIGH CRITICAL"`
	OutputIDs []*string `json:"outputIds"`
}

// SetDefaultOutputsOutput is the output of the SetDefaultOutputs operation
type SetDefaultOutputsOutput = DefaultOutputs

// GetDefaultOutputsInput is the request sent to return as part of GetDefaultOutputs operation
type GetDefaultOutputsInput struct {
}

// GetDefaultOutputsOutput is the response of the GetDefaultOutputs operation
type GetDefaultOutputsOutput struct {
	Defaults []*DefaultOutputs `json:"defaults"`
}

// AlertOutput contains the information for alert output configuration
type AlertOutput struct {

	// The user ID of the user that created the alert output
	CreatedBy *string `json:"createdBy"`

	// The time in epoch seconds when the alert output was created
	CreationTime *string `json:"creationTime"`

	// DisplayName is the user-provided name, e.g. "alert-channel".
	DisplayName *string `json:"displayName"`

	// The user ID of the user that last modified the alert output last
	LastModifiedBy *string `json:"lastModifiedBy"`

	// The time in epoch seconds when the alert output was last modified
	LastModifiedTime *string `json:"lastModifiedTime"`

	// Identifies uniquely an alert output (table sort key)
	OutputID *string `json:"outputId"`

	// OutputType is the output class, e.g. "slack", "sns".
	// ("type" is a reserved Dynamo keyword, so we use "OutputType" instead)
	OutputType *string `json:"outputType"`

	// OutputConfig contains the configuration for this output
	OutputConfig *OutputConfig `json:"outputConfig"`

	// VerificationStatus is the current state of the output verification process.
	VerificationStatus *string `json:"verificationStatus"`

	// DefaultForSeverity defines the alert severities that will be forwarded through this output
	DefaultForSeverity []*string `json:"defaultForSeverity"`
}

const (
	// VerificationStatusNotStarted shows that the verification process hasn't started yet
	VerificationStatusNotStarted = "NOT_STARTED"

	// VerificationStatusPending shows that a verification is pending
	VerificationStatusPending = "PENDING"

	// VerificationStatusFailed shows that the verification process has failed
	VerificationStatusFailed = "FAILED"

	// VerificationStatusSuccess shows that a verification is successful
	VerificationStatusSuccess = "SUCCESS"
)

// OutputConfig contains the configuration for the output
type OutputConfig struct {
	// SlackConfig contains the configuration for Slack alert output
	Slack *SlackConfig `json:"slack,omitempty"`

	// SnsConfig contains the configuration for SNS alert output
	Sns *SnsConfig `json:"sns,omitempty"`

	// SnsConfig contains the configuration for Email alert output
	Email *EmailConfig `json:"email,omitempty"`

	// PagerDuty contains the configuration for PagerDuty alert output
	PagerDuty *PagerDutyConfig `json:"pagerDuty,omitempty"`

	// Github contains the configuration for Github alert output
	Github *GithubConfig `json:"github,omitempty"`

	// Jira contains the configuration for Jira alert output
	Jira *JiraConfig `json:"jira,omitempty"`

	// Opsgenie contains the configuration for Opsgenie alert output
	Opsgenie *OpsgenieConfig `json:"opsgenie,omitempty"`

	// MsTeams contains the configuration for MsTeams alert output
	MsTeams *MsTeamsConfig `json:"msTeams,omitempty"`

	// SqsConfig contains the configuration for SQS alert output
	Sqs *SqsConfig `json:"sqs,omitempty"`
}

// SlackConfig defines options for each Slack output.
type SlackConfig struct {
	WebhookURL *string `json:"webhookURL" validate:"required,url"` // https://hooks.slack.com/services/...
}

// SnsConfig defines options for each SNS topic output
type SnsConfig struct {
	TopicArn *string `json:"topicArn" validate:"required,snsArn"`
}

// PagerDutyConfig defines options for each PagerDuty output
type PagerDutyConfig struct {
	IntegrationKey *string `json:"integrationKey" validate:"required,hexadecimal,len=32"`
}

// EmailConfig defines options for each Email output
type EmailConfig struct {
	DestinationAddress *string `json:"destinationAddress" validate:"required"`
}

// GithubConfig defines options for each Github output
type GithubConfig struct {
	RepoName *string `json:"repoName" validate:"required"`
	Token    *string `json:"token" validate:"required"`
}

// JiraConfig defines options for each Jira output
type JiraConfig struct {
	OrgDomain  *string `json:"orgDomain" validate:"required"`
	ProjectKey *string `json:"projectKey" validate:"required"`
	UserName   *string `json:"userName" validate:"required"`
	APIKey     *string `json:"apiKey" validate:"required"`
	AssigneeID *string `json:"assigneeId"`
}

// OpsgenieConfig defines options for each Opsgenie output
type OpsgenieConfig struct {
	APIKey *string `json:"apiKey" validate:"required"`
}

// MsTeamsConfig defines options for each MsTeamsConfig output
type MsTeamsConfig struct {
	WebhookURL *string `json:"webhookURL" validate:"required,url"`
}

// SqsConfig defines options for each Sqs topic output
type SqsConfig struct {
	QueueURL *string `json:"queueUrl" validate:"required,url"`
}

// DefaultOutputs is the structure holding the information about default outputs for severity
type DefaultOutputs struct {
	Severity  *string   `json:"severity"`
	OutputIDs []*string `json:"outputIds"`
}
