package aws

import "github.com/aws/aws-sdk-go/service/ec2"

const (
	Ec2VpcSchema = "AWS.EC2.VPC"
)

// Ec2Vpc contains all information about an EC2 VPC
type Ec2Vpc struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from ec2.Vpc
	CidrBlock                   *string
	CidrBlockAssociationSet     []*ec2.VpcCidrBlockAssociation
	DhcpOptionsId               *string
	InstanceTenancy             *string
	Ipv6CidrBlockAssociationSet []*ec2.VpcIpv6CidrBlockAssociation
	IsDefault                   *bool
	OwnerId                     *string
	State                       *string

	// Additional fields
	FlowLogs            []*ec2.FlowLog
	NetworkAcls         []*ec2.NetworkAcl
	RouteTables         []*ec2.RouteTable
	SecurityGroups      []*ec2.SecurityGroup
	StaleSecurityGroups []*ec2.StaleSecurityGroup
}
