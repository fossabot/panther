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

import "time"

// RuleType identifies the Alert to be for a Policy
const RuleType = "RULE"

// PolicyType identifies the Alert to be for a Policy
const PolicyType = "POLICY"

// Alert is the schema for each row in the Dynamo alerts table.
type Alert struct {

	// CreatedAt is the creation timestamp (seconds since epoch).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// OutputIDs is the set of outputs for this alert.
	OutputIDs []*string `json:"outputIds,omitempty"`

	// PolicyDescription is the description of the rule that triggered the alert.
	PolicyDescription *string `json:"policyDescription,omitempty"`

	// PolicyID is the rule that triggered the alert.
	PolicyID *string `json:"policyId" validate:"required"`

	// PolicyName is the name of the policy at the time the alert was triggered.
	PolicyName *string `json:"policyName,omitempty"`

	// PolicyVersionID is the S3 object version for the policy.
	PolicyVersionID *string `json:"policyVersionId,omitempty"`

	// Runbook is the user-provided triage information.
	Runbook *string `json:"runbook,omitempty"`

	// Severity is the alert severity at the time of creation.
	Severity *string `json:"severity" validate:"required,oneof=INFO LOW MEDIUM HIGH CRITICAL"`

	// Tags is the set of policy tags.
	Tags []*string `json:"tags,omitempty"`

	// AlertID specifies the alertId that this Alert is associated with.
	AlertID *string `json:"alertId,omitempty"`

	// Type specifies if an alert is for a policy or a rule
	Type *string `json:"type,omitempty" validate:"omitempty,oneof=RULE POLICY"`
}
