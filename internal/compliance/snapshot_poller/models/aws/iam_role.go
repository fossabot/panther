package aws

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

import "github.com/aws/aws-sdk-go/service/iam"

const (
	// IAMRoleSchema is the schema identifier for IAMRole.
	IAMRoleSchema = "AWS.IAM.Role"
)

// IAMRole contains all information about an IAM Role
type IAMRole struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from iam.Role
	AssumeRolePolicyDocument *string
	Description              *string
	MaxSessionDuration       *int64
	Path                     *string
	PermissionsBoundary      *iam.AttachedPermissionsBoundary

	// Additional fields
	InlinePolicies     map[string]*string
	ManagedPolicyNames []*string
}
