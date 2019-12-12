package aws

import "github.com/aws/aws-sdk-go/service/cloudtrail"

const (
	CloudTrailSchema     = "AWS.CloudTrail"
	CloudTrailMetaSchema = "AWS.CloudTrail.Meta"
)

// CloudTrail contains all information about a configured CloudTrail.
//
// This includes the trail info, status, event selectors, and attributes of the logging S3 bucket.
type CloudTrail struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from cloudtrail.Trail
	CloudWatchLogsLogGroupArn  *string
	CloudWatchLogsRoleArn      *string
	HasCustomEventSelectors    *bool
	HomeRegion                 *string
	IncludeGlobalServiceEvents *bool
	IsMultiRegionTrail         *bool
	IsOrganizationTrail        *bool
	KmsKeyId                   *string
	LogFileValidationEnabled   *bool
	S3BucketName               *string
	S3KeyPrefix                *string
	SnsTopicARN                *string
	SnsTopicName               *string // Deprecated by AWS

	// Additional fields
	EventSelectors []*cloudtrail.EventSelector
	Status         *cloudtrail.GetTrailStatusOutput
}

type CloudTrailMeta struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	Trails               []*string
	GlobalEventSelectors []*cloudtrail.EventSelector
}

// CloudTrails are a mapping of all Trails in an account keyed by ARN.
type CloudTrails map[string]*CloudTrail
