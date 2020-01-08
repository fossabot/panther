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
