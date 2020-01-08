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

import (
	"time"

	"github.com/aws/aws-sdk-go/service/rds"
)

const (
	RDSInstanceSchema = "AWS.RDS.Instance"
)

// RDSInstance contains all the information about an RDS DB instance
type RDSInstance struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from rds.DBInstance
	AllocatedStorage                      *int64
	AssociatedRoles                       []*rds.DBInstanceRole
	AutoMinorVersionUpgrade               *bool
	AvailabilityZone                      *string
	BackupRetentionPeriod                 *int64
	CACertificateIdentifier               *string
	CharacterSetName                      *string
	CopyTagsToSnapshot                    *bool
	DBClusterIdentifier                   *string
	DBInstanceClass                       *string
	DBInstanceStatus                      *string
	DBParameterGroups                     []*rds.DBParameterGroupStatus
	DBSecurityGroups                      []*rds.DBSecurityGroupMembership
	DBSubnetGroup                         *rds.DBSubnetGroup
	DbInstancePort                        *int64
	DbiResourceId                         *string
	DeletionProtection                    *bool
	DomainMemberships                     []*rds.DomainMembership
	EnabledCloudwatchLogsExports          []*string
	Endpoint                              *rds.Endpoint
	Engine                                *string
	EngineVersion                         *string
	EnhancedMonitoringResourceArn         *string
	IAMDatabaseAuthenticationEnabled      *bool
	Iops                                  *int64
	KmsKeyId                              *string
	LatestRestorableTime                  *time.Time
	LicenseModel                          *string
	ListenerEndpoint                      *rds.Endpoint
	MasterUsername                        *string
	MaxAllocatedStorage                   *int64
	MonitoringInterval                    *int64
	MonitoringRoleArn                     *string
	MultiAZ                               *bool
	OptionGroupMemberships                []*rds.OptionGroupMembership
	PendingModifiedValues                 *rds.PendingModifiedValues
	PerformanceInsightsEnabled            *bool
	PerformanceInsightsKMSKeyId           *string
	PerformanceInsightsRetentionPeriod    *int64
	PreferredBackupWindow                 *string
	PreferredMaintenanceWindow            *string
	ProcessorFeatures                     []*rds.ProcessorFeature
	PromotionTier                         *int64
	PubliclyAccessible                    *bool
	ReadReplicaDBClusterIdentifiers       []*string
	ReadReplicaDBInstanceIdentifiers      []*string
	ReadReplicaSourceDBInstanceIdentifier *string
	SecondaryAvailabilityZone             *string
	StatusInfos                           []*rds.DBInstanceStatusInfo
	StorageEncrypted                      *bool
	StorageType                           *string
	TdeCredentialArn                      *string
	Timezone                              *string
	VpcSecurityGroups                     []*rds.VpcSecurityGroupMembership

	// Additional fields
	SnapshotAttributes []*rds.DBSnapshotAttributesResult
}
