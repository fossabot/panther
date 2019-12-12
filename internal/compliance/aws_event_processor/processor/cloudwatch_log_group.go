package processor

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyCloudWatchLogGroup(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_amazoncloudwatchlogs.html
	if eventName == "CancelExportTask" ||
		eventName == "CreateExportTask" ||
		eventName == "PutDestination" ||
		eventName == "PutDestinationPolicy" ||
		eventName == "PutLogEvents" ||
		eventName == "PutResourcePolicy" ||
		eventName == "StartQuery" ||
		eventName == "StopQuery" ||
		eventName == "TestMetricFilter" ||
		eventName == "CreateLogStream" ||
		eventName == "FilterLogEvents" {

		zap.L().Debug("loggroup: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	region := detail.Get("awsRegion").Str
	logGroupARN := arn.ARN{
		Partition: "aws",
		Service:   "logs",
		Region:    region,
		AccountID: accountID,
		Resource:  "log-group:",
	}
	switch eventName {
	case "AssociateKmsKey", "CreateLogGroup", "DeleteLogGroup", "DeleteLogStream", "DeleteMetricFilter",
		"DeleteRetentionPolicy", "DeleteSubscriptionFilter", "DisassociateKmsKey", "PutMetricFilter",
		"PutRetentionPolicy", "PutSubscriptionFilter", "TagLogGroup", "UntagLogGroup":
		// Not technically the correct resourceID, see classifyCloudFormation for a more detailed
		// explanation.
		logGroupARN.Resource += detail.Get("requestParameters.logGroupName").Str
	default:
		zap.L().Warn("loggroup: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteLogGroup",
		EventName:    eventName,
		ResourceID:   logGroupARN.String(),
		ResourceType: schemas.CloudWatchLogGroupSchema,
	}}
}
