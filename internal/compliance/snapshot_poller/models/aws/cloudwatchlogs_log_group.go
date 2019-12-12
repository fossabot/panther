package aws

const (
	CloudWatchLogGroupSchema = "AWS.CloudWatch.LogGroup"
)

// CloudWatchLogsLogGroup contains all the information about an CloudWatch Logs Log Group
type CloudWatchLogsLogGroup struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from cloudwatchlogs.LogGroup
	KmsKeyId          *string
	MetricFilterCount *int64
	RetentionInDays   *int64
	StoredBytes       *int64
}
