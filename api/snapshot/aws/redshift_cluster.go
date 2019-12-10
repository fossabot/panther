package aws

import "github.com/aws/aws-sdk-go/service/redshift"

const (
	RedshiftClusterSchema = "AWS.Redshift.Cluster"
)

// RedshiftCluseter contains all the information about a Redshift cluster
type RedshiftCluster struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from redshift.cluster
	AllowVersionUpgrade              *bool
	AutomatedSnapshotRetentionPeriod *int64
	AvailabilityZone                 *string
	ClusterAvailabilityStatus        *string
	ClusterNodes                     []*redshift.ClusterNode
	ClusterParameterGroups           []*redshift.ClusterParameterGroupStatus
	ClusterPublicKey                 *string
	ClusterRevisionNumber            *string
	ClusterSecurityGroups            []*redshift.ClusterSecurityGroupMembership
	ClusterSnapshotCopyStatus        *redshift.ClusterSnapshotCopyStatus
	ClusterStatus                    *string
	ClusterSubnetGroupName           *string
	ClusterVersion                   *string
	DataTransferProgress             *redshift.DataTransferProgress
	DeferredMaintenanceWindows       []*redshift.DeferredMaintenanceWindow
	ElasticIpStatus                  *redshift.ElasticIpStatus
	ElasticResizeNumberOfNodeOptions *string
	Encrypted                        *bool
	Endpoint                         *redshift.Endpoint
	EnhancedVpcRouting               *bool
	HsmStatus                        *redshift.HsmStatus
	IamRoles                         []*redshift.ClusterIamRole
	KmsKeyId                         *string
	MaintenanceTrackName             *string
	ManualSnapshotRetentionPeriod    *int64
	MasterUsername                   *string
	ModifyStatus                     *string
	NodeType                         *string
	NumberOfNodes                    *int64
	PendingActions                   []*string
	PendingModifiedValues            *redshift.PendingModifiedValues
	PreferredMaintenanceWindow       *string
	PubliclyAccessible               *bool
	ResizeInfo                       *redshift.ResizeInfo
	RestoreStatus                    *redshift.RestoreStatus
	SnapshotScheduleIdentifier       *string
	SnapshotScheduleState            *string
	VpcId                            *string
	VpcSecurityGroups                []*redshift.VpcSecurityGroupMembership

	// Additional fields
	LoggingStatus *redshift.LoggingStatus
}
