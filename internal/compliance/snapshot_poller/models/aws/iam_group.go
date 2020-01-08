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
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	IAMGroupSchema = "AWS.IAM.Group"
)

// IamGroup contains all the information about an IAM Group
type IamGroup struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from iam.Group
	Path *string

	// Additional fields
	InlinePolicies    map[string]*string
	ManagedPolicyARNs []*string
	Users             []*iam.User
}
