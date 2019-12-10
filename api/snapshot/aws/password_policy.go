package aws

import "github.com/aws/aws-sdk-go/service/iam"

const (
	PasswordPolicySchema = "AWS.PasswordPolicy"
)

// PasswordPolicy contains all information about a configured password policy.
type PasswordPolicy struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	iam.PasswordPolicy
	AnyExist bool
}
