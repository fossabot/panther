package aws

import "github.com/aws/aws-sdk-go/service/ec2"

const (
	Ec2SecurityGroupSchema = "AWS.EC2.SecurityGroup"
)

// Ec2SecurityGroup contains all information about an EC2 SecurityGroup
type Ec2SecurityGroup struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from ec2.SecurityGroup
	Description         *string
	IpPermissions       []*ec2.IpPermission
	IpPermissionsEgress []*ec2.IpPermission
	OwnerId             *string
	VpcId               *string
}
