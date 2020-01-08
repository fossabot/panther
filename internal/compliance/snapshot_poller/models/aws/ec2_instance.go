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
