package aws

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
