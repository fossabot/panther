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

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	CloudFormationStackSchema = "AWS.CloudFormation.Stack"
)

// CloudFormationStack contains all the information about a CloudFormation Stack
type CloudFormationStack struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from cloudformation.Stack
	Capabilities                []*string
	ChangeSetId                 *string
	DeletionTime                *time.Time
	Description                 *string
	DisableRollback             *bool
	DriftInformation            *cloudformation.StackDriftInformation
	EnableTerminationProtection *bool
	LastUpdatedTime             *time.Time
	NotificationARNs            []*string
	Outputs                     []*cloudformation.Output
	Parameters                  []*cloudformation.Parameter
	ParentId                    *string
	RoleARN                     *string
	RollbackConfiguration       *cloudformation.RollbackConfiguration
	RootId                      *string
	StackStatus                 *string
	StackStatusReason           *string
	TimeoutInMinutes            *int64

	// Additional fields
	Drifts []*cloudformation.StackResourceDrift
}
