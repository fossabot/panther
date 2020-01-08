package aws

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
