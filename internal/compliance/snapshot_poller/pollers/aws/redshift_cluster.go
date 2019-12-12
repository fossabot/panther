package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshift/redshiftiface"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/gateway/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

// Set as variables to be overridden in testing
var (
	RedshiftClientFunc = setupRedshiftClient
)

func setupRedshiftClient(sess *session.Session, cfg *aws.Config) interface{} {
	cfg.MaxRetries = aws.Int(MaxRetries)
	return redshift.New(sess, cfg)
}

// PollRedshiftCluster polls a single Redshift Cluster resource
func PollRedshiftCluster(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	scanRequest *pollermodels.ScanEntry,
) interface{} {

	client := getClient(pollerResourceInput, "redshift", resourceARN.Region).(redshiftiface.RedshiftAPI)
	clusterID := strings.Replace(resourceARN.Resource, "cluster:", "", 1)
	redshiftCluster := getRedshiftCluster(client, aws.String(clusterID))

	snapshot := buildRedshiftClusterSnapshot(client, redshiftCluster)
	if snapshot == nil {
		return nil
	}
	snapshot.ResourceID = scanRequest.ResourceID
	snapshot.AccountID = aws.String(resourceARN.AccountID)
	snapshot.Region = aws.String(resourceARN.Region)
	snapshot.ARN = scanRequest.ResourceID
	return snapshot
}

// getRedshiftCluster returns a specific redshift cluster
func getRedshiftCluster(svc redshiftiface.RedshiftAPI, clusterID *string) *redshift.Cluster {
	cluster, err := svc.DescribeClusters(&redshift.DescribeClustersInput{
		ClusterIdentifier: clusterID,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ClusterNotFound" {
				zap.L().Warn("tried to scan non-existent resource",
					zap.String("resource", *clusterID),
					zap.String("resourceType", awsmodels.RedshiftClusterSchema))
				return nil
			}
		}
		utils.LogAWSError("Redshift.DescribeClusters", err)
		return nil
	}
	return cluster.Clusters[0]
}

// describeClusters returns a list of all redshift cluster in the account
func describeClusters(redshiftSvc redshiftiface.RedshiftAPI) (clusters []*redshift.Cluster) {
	err := redshiftSvc.DescribeClustersPages(&redshift.DescribeClustersInput{},
		func(page *redshift.DescribeClustersOutput, lastPage bool) bool {
			clusters = append(clusters, page.Clusters...)
			return true
		})
	if err != nil {
		utils.LogAWSError("Redshift.DescribeClustersPages", err)
	}
	return
}

// describeLoggingStatus determines whether or not a redshift cluster has logging enabled
func describeLoggingStatus(redshiftSvc redshiftiface.RedshiftAPI, clusterID *string) (*redshift.LoggingStatus, error) {
	out, err := redshiftSvc.DescribeLoggingStatus(
		&redshift.DescribeLoggingStatusInput{ClusterIdentifier: clusterID},
	)
	if err != nil {
		utils.LogAWSError("Redshift.DescribeLoggingStatus", err)
		return nil, err
	}
	return out, nil
}

// buildRedshiftClusterSnapshot makes all the calls to build up a snapshot of a given Redshift cluster
func buildRedshiftClusterSnapshot(redshiftSvc redshiftiface.RedshiftAPI, cluster *redshift.Cluster) *awsmodels.RedshiftCluster {
	if cluster == nil {
		return nil
	}
	clusterSnapshot := &awsmodels.RedshiftCluster{
		GenericResource: awsmodels.GenericResource{
			TimeCreated:  utils.DateTimeFormat(*cluster.ClusterCreateTime),
			ResourceType: aws.String(awsmodels.RedshiftClusterSchema),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			Name: cluster.DBName,
			ID:   cluster.ClusterIdentifier,
			Tags: utils.ParseTagSlice(cluster.Tags),
		},
		AllowVersionUpgrade:              cluster.AllowVersionUpgrade,
		AutomatedSnapshotRetentionPeriod: cluster.AutomatedSnapshotRetentionPeriod,
		AvailabilityZone:                 cluster.AvailabilityZone,
		ClusterAvailabilityStatus:        cluster.ClusterAvailabilityStatus,
		ClusterNodes:                     cluster.ClusterNodes,
		ClusterParameterGroups:           cluster.ClusterParameterGroups,
		ClusterPublicKey:                 cluster.ClusterPublicKey,
		ClusterRevisionNumber:            cluster.ClusterRevisionNumber,
		ClusterSecurityGroups:            cluster.ClusterSecurityGroups,
		ClusterSnapshotCopyStatus:        cluster.ClusterSnapshotCopyStatus,
		ClusterStatus:                    cluster.ClusterStatus,
		ClusterSubnetGroupName:           cluster.ClusterSubnetGroupName,
		ClusterVersion:                   cluster.ClusterVersion,
		DataTransferProgress:             cluster.DataTransferProgress,
		DeferredMaintenanceWindows:       cluster.DeferredMaintenanceWindows,
		ElasticIpStatus:                  cluster.ElasticIpStatus,
		ElasticResizeNumberOfNodeOptions: cluster.ElasticResizeNumberOfNodeOptions,
		Encrypted:                        cluster.Encrypted,
		Endpoint:                         cluster.Endpoint,
		EnhancedVpcRouting:               cluster.EnhancedVpcRouting,
		HsmStatus:                        cluster.HsmStatus,
		IamRoles:                         cluster.IamRoles,
		KmsKeyId:                         cluster.KmsKeyId,
		MaintenanceTrackName:             cluster.MaintenanceTrackName,
		ManualSnapshotRetentionPeriod:    cluster.ManualSnapshotRetentionPeriod,
		MasterUsername:                   cluster.MasterUsername,
		ModifyStatus:                     cluster.ModifyStatus,
		NodeType:                         cluster.NodeType,
		NumberOfNodes:                    cluster.NumberOfNodes,
		PendingActions:                   cluster.PendingActions,
		PendingModifiedValues:            cluster.PendingModifiedValues,
		PreferredMaintenanceWindow:       cluster.PreferredMaintenanceWindow,
		PubliclyAccessible:               cluster.PubliclyAccessible,
		ResizeInfo:                       cluster.ResizeInfo,
		RestoreStatus:                    cluster.RestoreStatus,
		SnapshotScheduleIdentifier:       cluster.SnapshotScheduleIdentifier,
		SnapshotScheduleState:            cluster.SnapshotScheduleState,
		VpcId:                            cluster.VpcId,
		VpcSecurityGroups:                cluster.VpcSecurityGroups,
	}

	loggingStatus, err := describeLoggingStatus(redshiftSvc, cluster.ClusterIdentifier)
	if err != nil {
		return clusterSnapshot
	}
	clusterSnapshot.LoggingStatus = loggingStatus

	return clusterSnapshot
}

// PollRedshiftClusters gathers information on each Redshift Cluster for an AWS account.
func PollRedshiftClusters(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting Redshift Cluster resource poller")
	redshiftClusterSnapshots := make(map[string]*awsmodels.RedshiftCluster)

	for _, regionID := range utils.GetServiceRegions(pollerInput.Regions, "redshift") {
		sess := session.Must(session.NewSession(&aws.Config{Region: regionID}))
		creds, err := AssumeRoleFunc(pollerInput, sess)
		if err != nil {
			return nil, err
		}

		redshiftSvc := RedshiftClientFunc(sess, &aws.Config{Credentials: creds}).(redshiftiface.RedshiftAPI)

		// Start with generating a list of all clusters
		clusters := describeClusters(redshiftSvc)
		if len(clusters) == 0 {
			zap.L().Debug("No Redshift clusters found.", zap.String("region", *regionID))
			continue
		}

		for _, cluster := range clusters {
			redshiftClusterSnapshot := buildRedshiftClusterSnapshot(redshiftSvc, cluster)

			resourceID := strings.Join(
				[]string{
					"arn",
					pollerInput.AuthSourceParsedARN.Partition,
					"redshift",
					*regionID,
					pollerInput.AuthSourceParsedARN.AccountID,
					"cluster",
					*redshiftClusterSnapshot.ID},
				":",
			)
			// Populate generic fields
			redshiftClusterSnapshot.ResourceID = aws.String(resourceID)

			// Populate AWS generic fields
			redshiftClusterSnapshot.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)
			redshiftClusterSnapshot.Region = regionID
			redshiftClusterSnapshot.ARN = aws.String(resourceID)

			if _, ok := redshiftClusterSnapshots[resourceID]; !ok {
				redshiftClusterSnapshots[resourceID] = redshiftClusterSnapshot
			} else {
				zap.L().Info(
					"overwriting existing Redshift Cluster snapshot",
					zap.String("resourceID", resourceID),
				)
				redshiftClusterSnapshots[resourceID] = redshiftClusterSnapshot
			}
		}
	}

	resources := make([]*apimodels.AddResourceEntry, 0, len(redshiftClusterSnapshots))
	for resourceID, clusterSnapshot := range redshiftClusterSnapshots {
		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      clusterSnapshot,
			ID:              apimodels.ResourceID(resourceID),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.RedshiftClusterSchema,
		})
	}

	return resources, nil
}