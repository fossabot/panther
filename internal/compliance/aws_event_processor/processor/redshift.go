package processor

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyRedshift(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	if eventName == "AcceptReservedNodeExchange" ||
		eventName == "CreateClusterSecurityGroup" ||
		eventName == "CreateHsmClientCertificate" ||
		eventName == "CreateHsmConfiguration" ||
		eventName == "DeleteClusterParameterGroup" ||
		eventName == "DeleteClusterSecurityGroup" ||
		eventName == "DeleteClusterSubnetGroup" ||
		eventName == "DeleteEventSubscription" ||
		eventName == "DeleteHsmClientCertificate" ||
		eventName == "DeleteHsmConfiguration" ||
		eventName == "DeleteSnapshotCopyGrant" ||
		eventName == "DeleteSnapshotSchedule" ||
		eventName == "ModifyClusterParameterGroup" ||
		eventName == "ModifyClusterSubnetGroup" ||
		eventName == "ResetClusterParameterGroup" ||
		eventName == "RevokeClusterSecurityGroupIngress" ||
		eventName == "CreateClusterParameterGroup" {

		zap.L().Debug("redshift: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_amazonredshift.html
	region := detail.Get("awsRegion").Str
	redshiftARN := arn.ARN{
		Partition: "aws",
		Service:   "redshift",
		Region:    region,
		AccountID: accountID,
		Resource:  "cluster:",
	}

	// CreateClusterUser???
	switch eventName {
	case "AuthorizeSnapshotAccess", "CopyClusterSnapshot", "DeleteClusterSnapshot", "ModifyClusterSnapshot", "RevokeSnapshotAccess":
		// If we add a cluster snapshot resource, this should be updated to include that as well
		redshiftARN.Resource += detail.Get("responseElements.snapshot.clusterIdentifier").Str
	case "AuthorizeClusterSecurityGroupIngress", "BatchDeleteClusterSnapshots", "BatchModifyClusterSnapshots", "CreateSnapshotCopyGrant":
		// We don't have a good way to tie this back to a single cluster, so do a region wide scan
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			Region:       region,
			ResourceType: schemas.RedshiftClusterSchema,
		}}
	case "CancelResize", "CreateCluster", "CreateClusterSnapshot", "DeleteCluster", "DisableLogging", "DisableSnapshotCopy",
		"EnableLogging", "EnableSnapshotCopy", "ModifyCluster", "ModifyClusterDbRevision", "ModifyClusterIamRoles",
		"ModifyClusterMaintenance", "ModifyClusterSnapshotSchedule", "ModifySnapshotCopyRetentionPeriod", "RebootCluster",
		"ResizeCluster", "RestoreFromClusterSnapshot", "ClusterIdentifier", "RotateEncryptionKey":
		redshiftARN.Resource += detail.Get("requestParameters.clusterIdentifier").Str
	case "CreateClusterSubnetGroup":
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			ResourceID: arn.ARN{
				Partition: "aws",
				Service:   "ec2",
				Region:    region,
				AccountID: accountID,
				Resource:  "vpc/" + detail.Get("responseElements.clusterSubnetGroup.vpcId").Str,
			}.String(),
			ResourceType: schemas.Ec2VpcSchema,
		}}
	case "CreateEventSubscription", "ModifyEventSubscription":
		// Optional parameters
		resourceType := detail.Get("requestParameters.sourceType").Str
		resourceIDs := detail.Get("requestParameters.sourceIds").Array()

		// Handle the case where this is explicitly for non-cluster resources
		if resourceType != "cluster" && resourceType != "" {
			return nil
		}

		// Handle the cases where this is for all resources, or all cluster resources
		if resourceType == "" || len(resourceIDs) == 0 {
			// The documentation says this means it applies to every Redshift resource in the
			// account, but I suspect it may be regional. Further testing required.
			return []*resourceChange{{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceType: schemas.RedshiftClusterSchema,
			}}
		}

		// Handle the case where this is for an explicit list of resources
		changes := make([]*resourceChange, len(resourceIDs))
		for i, resourceID := range resourceIDs {
			changes[i] = &resourceChange{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceID:   redshiftARN.String() + resourceID.Str,
				ResourceType: schemas.RedshiftClusterSchema,
			}
		}
		return changes
	case "CreateSnapshotSchedule", "ModifySnapshotSchedule":
		clusters := detail.Get("responseElements.associatedClusters").Array()
		changes := make([]*resourceChange, len(clusters))
		for i, cluster := range clusters {
			changes[i] = &resourceChange{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceID:   redshiftARN.String() + cluster.Get("clusterIdentifier").Str,
				ResourceType: schemas.RedshiftClusterSchema,
			}
		}
		return changes
	case "CreateTags", "DeleteTags":
		resourceARN, err := arn.Parse(detail.Get("requestParameters.resourceName").Str)
		if err != nil {
			zap.L().Error("redshift: error parsing ARN", zap.String("eventName", eventName), zap.Error(err))
			return nil
		}
		if strings.HasPrefix(resourceARN.Resource, "cluster:") {
			redshiftARN = resourceARN
			break
		}
		return nil
	default:
		zap.L().Warn("redshift: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteCluster",
		EventName:    eventName,
		ResourceID:   redshiftARN.String(),
		ResourceType: schemas.RedshiftClusterSchema,
	}}
}
