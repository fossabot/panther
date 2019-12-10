package aws

import "github.com/aws/aws-sdk-go/service/ec2"

const (
	Ec2AmiSchema = "AWS.EC2.AMI"
)

// Ec2Ami contains all information about an EC2 AMI
type Ec2Ami struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from ec2.Image
	Architecture        *string
	BlockDeviceMappings []*ec2.BlockDeviceMapping
	Description         *string
	EnaSupport          *bool
	Hypervisor          *string
	ImageLocation       *string
	ImageOwnerAlias     *string
	ImageType           *string
	KernelId            *string
	OwnerId             *string
	Platform            *string
	ProductCodes        []*ec2.ProductCode
	Public              *bool
	RamdiskId           *string
	RootDeviceName      *string
	RootDeviceType      *string
	SriovNetSupport     *string
	State               *string
	StateReason         *ec2.StateReason
	VirtualizationType  *string
}
