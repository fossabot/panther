package aws

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
