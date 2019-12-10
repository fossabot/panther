package processor

import (
	"strings"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// CloudWatch events which require downstream processing are summarized with this struct.
type resourceChange struct {
	AwsAccountID  string `json:"awsAccountId"`  // the 12-digit AWS account ID which owns the resource
	Delay         int64  `json:"delay"`         // How long in seconds to delay this message in SQS
	Delete        bool   `json:"delete"`        // True if the resource should be marked deleted (otherwise, update)
	EventName     string `json:"eventName"`     // CloudTrail event name (for logging only)
	EventTime     string `json:"eventTime"`     // official CloudTrail RFC3339 timestamp
	IntegrationID string `json:"integrationId"` // account integration ID
	Region        string `json:"region"`        // Region (for resource type scans only)
	ResourceID    string `json:"resourceId"`    // e.g. "arn:aws:s3:::my-bucket"
	ResourceType  string `json:"resourceType"`  // e.g. "AWS.S3.Bucket"
}

// Map each event source to the appropriate classifier function.
//
// The "classifier" takes a cloudtrail log and summarizes the required change.
// integrationID does not need to be set by the individual classifiers.
var classifiers = map[string]func(gjson.Result, string) []*resourceChange{
	"aws.acm":                  classifyACM,
	"aws.cloudformation":       classifyCloudFormation,
	"aws.cloudtrail":           classifyCloudTrail,
	"aws.config":               classifyConfig,
	"aws.dynamodb":             classifyDynamoDB,
	"aws.ec2":                  classifyEC2,
	"aws.elasticloadbalancing": classifyELBV2,
	"aws.guardduty":            classifyGuardDuty,
	"aws.iam":                  classifyIAM,
	"aws.kms":                  classifyKMS,
	"aws.lambda":               classifyLambda,
	"aws.logs":                 classifyCloudWatchLogGroup,
	"aws.rds":                  classifyRDS,
	"aws.redshift":             classifyRedshift,
	"aws.s3":                   classifyS3,
	"aws.waf":                  classifyWAF,
	"aws.waf-regional":         classifyWAFRegional,
}

// Classify the event as an update or delete operation, or drop it altogether.
func classifyCloudWatchEvent(body string) []*resourceChange {
	source := gjson.Get(body, "source").Str
	classifier, ok := classifiers[source]
	if !ok {
		zap.L().Debug("dropping event from unsupported source", zap.String("eventSource", source))
		return nil
	}

	// NOTE: we ignore the "detail.readOnly" field because it is not consistent
	detail := gjson.Get(body, "detail")
	if errorCode := detail.Get("errorCode").Str; errorCode != "" {
		zap.L().Debug("dropping failed event",
			zap.String("eventSource", source),
			zap.String("errorCode", errorCode))
		return nil
	}

	// Ignore the most common read only resources
	eventName := detail.Get("eventName").Str
	if strings.HasPrefix(eventName, "Get") ||
		strings.HasPrefix(eventName, "BatchGet") ||
		strings.HasPrefix(eventName, "Describe") ||
		strings.HasPrefix(eventName, "List") {

		zap.L().Debug(source+": ignoring read-only event", zap.String("eventName", eventName))
		return nil
	}

	// Check if this log is from a supported account
	accountID := gjson.Get(body, "account").Str
	integration, ok := accounts[accountID]
	if !ok {
		zap.L().Warn("dropping event from unauthorized account",
			zap.String("accountId", accountID),
			zap.String("eventSource", source))
		return nil
	}

	// Process the body
	changes := classifier(detail, accountID)
	eventTime := detail.Get("eventTime").Str
	for _, change := range changes {
		change.EventTime = eventTime
		change.IntegrationID = *integration.IntegrationID
	}

	return changes
}
