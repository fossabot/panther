package aws

import "github.com/aws/aws-sdk-go/service/ec2"

const (
	Ec2InstanceSchema = "AWS.EC2.Instance"
)

// Ec2Instance contains all information about an EC2 Instance
type Ec2Instance struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from ec2.Instance
	AmiLaunchIndex                          *int64
	Architecture                            *string
	BlockDeviceMappings                     []*ec2.InstanceBlockDeviceMapping
	CapacityReservationId                   *string
	CapacityReservationSpecification        *ec2.CapacityReservationSpecificationResponse
	ClientToken                             *string
	CpuOptions                              *ec2.CpuOptions
	EbsOptimized                            *bool
	ElasticGpuAssociations                  []*ec2.ElasticGpuAssociation
	ElasticInferenceAcceleratorAssociations []*ec2.ElasticInferenceAcceleratorAssociation
	EnaSupport                              *bool
	HibernationOptions                      *ec2.HibernationOptions
	Hypervisor                              *string
	IamInstanceProfile                      *ec2.IamInstanceProfile
	ImageId                                 *string
	InstanceLifecycle                       *string
	InstanceType                            *string
	KernelId                                *string
	KeyName                                 *string
	Licenses                                []*ec2.LicenseConfiguration
	Monitoring                              *ec2.Monitoring
	NetworkInterfaces                       []*ec2.InstanceNetworkInterface
	Placement                               *ec2.Placement
	Platform                                *string
	PrivateDnsName                          *string
	PrivateIpAddress                        *string
	ProductCodes                            []*ec2.ProductCode
	PublicDnsName                           *string
	PublicIpAddress                         *string
	RamdiskId                               *string
	RootDeviceName                          *string
	RootDeviceType                          *string
	SecurityGroups                          []*ec2.GroupIdentifier
	SourceDestCheck                         *bool
	SpotInstanceRequestId                   *string
	SriovNetSupport                         *string
	State                                   *ec2.InstanceState
	StateReason                             *ec2.StateReason
	StateTransitionReason                   *string
	SubnetId                                *string
	VirtualizationType                      *string
	VpcId                                   *string
}
