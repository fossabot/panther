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
