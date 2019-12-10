package aws

import "github.com/aws/aws-sdk-go/service/ec2"

const (
	Ec2NetworkAclSchema = "AWS.EC2.NetworkACL"
)

// Ec2NetworkACL contains all information about an EC2 Network ACL
type Ec2NetworkAcl struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from ec2.NetworkAcl
	Associations []*ec2.NetworkAclAssociation
	Entries      []*ec2.NetworkAclEntry
	IsDefault    *bool
	OwnerId      *string
	VpcId        *string
}
