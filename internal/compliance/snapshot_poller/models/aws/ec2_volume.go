package aws

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
