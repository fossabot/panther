package aws

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
