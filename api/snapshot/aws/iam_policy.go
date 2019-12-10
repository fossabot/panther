package aws

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
