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
	Ec2VolumeSchema = "AWS.EC2.Volume"
)

// Ec2Volume contains all the information about an EC2 Volume
type Ec2Volume struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from ec2.Volume
	Attachments      []*ec2.VolumeAttachment
	AvailabilityZone *string
	Encrypted        *bool
	Iops             *int64
	KmsKeyId         *string
	Size             *int64
	SnapshotId       *string
	State            *string
	VolumeType       *string

	// Additional fields
	Snapshots []*Ec2Snapshot
}

type Ec2Snapshot struct {
	*ec2.Snapshot
	CreateVolumePermissions []*ec2.CreateVolumePermission
}
