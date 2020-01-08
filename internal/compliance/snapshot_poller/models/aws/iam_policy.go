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

import (
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	IAMPolicySchema = "AWS.IAM.Policy"
)

// IAMPolicy contains all information about a policy.
type IAMPolicy struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from iam.Policy
	AttachmentCount               *int64
	DefaultVersionId              *string
	Description                   *string
	IsAttachable                  *bool
	Path                          *string
	PermissionsBoundaryUsageCount *int64
	UpdateDate                    *time.Time

	// Additional fields
	Entities       *IAMPolicyEntities
	PolicyDocument *string
}

// IAMPolicyEntities provides detail on the attached entities to an IAM policy.
type IAMPolicyEntities struct {
	PolicyGroups []*iam.PolicyGroup
	PolicyRoles  []*iam.PolicyRole
	PolicyUsers  []*iam.PolicyUser
}
