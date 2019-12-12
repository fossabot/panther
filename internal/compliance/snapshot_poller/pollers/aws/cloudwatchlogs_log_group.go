package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/gateway/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

// Set as variables to be overridden in testing
var (
	CloudWatchLogsClientFunc = setupCloudWatchLogsClient
)

func setupCloudWatchLogsClient(sess *session.Session, cfg *aws.Config) interface{} {
	cfg.MaxRetries = aws.Int(MaxRetries)
	return cloudwatchlogs.New(sess, cfg)
}

// PollCloudWatchLogsLogGroup polls a single CloudWatchLogs LogGroup resource
func PollCloudWatchLogsLogGroup(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	scanRequest *pollermodels.ScanEntry) (resource interface{}) {

	client := getClient(pollerResourceInput, "cloudwatchlogs", resourceARN.Region).(cloudwatchlogsiface.CloudWatchLogsAPI)

	// See PollCloudFormationStack for a detailed reasoning behind these actions
	// Get just the resource portion of the ARN, drop the resource type prefix
	lgResource := strings.TrimPrefix(resourceARN.Resource, "log-group:")

	// Split out the log group name from any additional modifiers
	lgName := strings.Split(lgResource, ":")[0]
	logGroup := getLogGroup(client, aws.String(lgName))
	snapshot := buildCloudWatchLogsLogGroupSnapshot(client, logGroup)
	if snapshot == nil {
		return nil
	}
	snapshot.Region = aws.String(resourceARN.Region)
	snapshot.AccountID = aws.String(resourceARN.AccountID)
	scanRequest.ResourceID = snapshot.ARN
	return snapshot
}

// getLogGroup returns a specific cloudwatch logs log group
func getLogGroup(svc cloudwatchlogsiface.CloudWatchLogsAPI, logGroupName *string) *cloudwatchlogs.LogGroup {
	logGroups, err := svc.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: logGroupName,
	})
	if err != nil {
		utils.LogAWSError("CloudWatchLogs.DescribeLogGroups", err)
		return nil
	}

	for _, logGroup := range logGroups.LogGroups {
		if *logGroup.LogGroupName == *logGroupName {
			return logGroup
		}
	}

	zap.L().Warn("tried to scan non-existent resource",
		zap.String("resource", *logGroupName),
		zap.String("resourceType", awsmodels.CloudWatchLogGroupSchema))
	return nil
}

// describeLogGroups returns all Log Groups in the account
func describeLogGroups(cloudwatchLogsSvc cloudwatchlogsiface.CloudWatchLogsAPI) (logGroups []*cloudwatchlogs.LogGroup) {
	err := cloudwatchLogsSvc.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{},
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			logGroups = append(logGroups, page.LogGroups...)
			return true
		})
	if err != nil {
		utils.LogAWSError("CloudWatchLogs.DescribeLogGroups", err)
	}
	return
}

// listTagsLogGroup returns the tags for a given log group
func listTagsLogGroup(svc cloudwatchlogsiface.CloudWatchLogsAPI, groupName *string) map[string]*string {
	tags, err := svc.ListTagsLogGroup(&cloudwatchlogs.ListTagsLogGroupInput{
		LogGroupName: groupName,
	})
	if err != nil {
		utils.LogAWSError("CloudWatchLogs ListTagsLogGroup", err)
		return nil
	}
	return tags.Tags
}

// buildCloudWatchLogsLogGroupSnapshot returns a complete snapshot of a LogGroup
func buildCloudWatchLogsLogGroupSnapshot(
	svc cloudwatchlogsiface.CloudWatchLogsAPI,
	logGroup *cloudwatchlogs.LogGroup,
) *awsmodels.CloudWatchLogsLogGroup {

	logGroupSnapshot := &awsmodels.CloudWatchLogsLogGroup{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   logGroup.Arn,
			ResourceType: aws.String(awsmodels.CloudWatchLogGroupSchema),
			// Convert milliseconds to seconds before converting to datetime
			// loses nanosecond precision
			TimeCreated: utils.UnixTimeToDateTime(*logGroup.CreationTime / 1000),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			Name: logGroup.LogGroupName,
			ARN:  logGroup.Arn,
		},
		KmsKeyId:          logGroup.KmsKeyId,
		MetricFilterCount: logGroup.MetricFilterCount,
		RetentionInDays:   logGroup.RetentionInDays,
		StoredBytes:       logGroup.StoredBytes,
	}
	logGroupSnapshot.Tags = listTagsLogGroup(svc, logGroupSnapshot.Name)

	return logGroupSnapshot
}

// PollCloudWatchLogsLogGroups gathers information on each CloudWatchLogs LogGroup for an AWS account
func PollCloudWatchLogsLogGroups(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting CloudWatch LogGroup resource poller")
	logGroupSnapshots := make(map[string]*awsmodels.CloudWatchLogsLogGroup)

	for _, regionID := range utils.GetServiceRegions(pollerInput.Regions, "logs") {
		sess := session.Must(session.NewSession(&aws.Config{Region: regionID}))
		creds, err := AssumeRoleFunc(pollerInput, sess)
		if err != nil {
			return nil, err
		}

		var cloudwatchLogGroupSvc = CloudWatchLogsClientFunc(sess, &aws.Config{Credentials: creds}).(cloudwatchlogsiface.CloudWatchLogsAPI)

		// Start with generating a list of all log groups
		logGroups := describeLogGroups(cloudwatchLogGroupSvc)
		if len(logGroups) == 0 {
			zap.L().Debug("no CloudWatchLogs LogGroups found", zap.String("region", *regionID))
			continue
		}

		for _, logGroup := range logGroups {
			logGroupSnapshot := buildCloudWatchLogsLogGroupSnapshot(cloudwatchLogGroupSvc, logGroup)
			if logGroupSnapshot == nil {
				continue
			}
			logGroupSnapshot.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)
			logGroupSnapshot.Region = regionID

			if _, ok := logGroupSnapshots[*logGroupSnapshot.ARN]; !ok {
				logGroupSnapshots[*logGroup.Arn] = logGroupSnapshot
			} else {
				zap.L().Info(
					"overwriting existing CloudWatchLogs LogGroup snapshot",
					zap.String("resourceId", *logGroupSnapshot.ARN),
				)
				logGroupSnapshots[*logGroupSnapshot.ARN] = logGroupSnapshot
			}
		}
	}

	resources := make([]*apimodels.AddResourceEntry, 0, len(logGroupSnapshots))
	for resourceID, logGroupSnapshot := range logGroupSnapshots {
		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      logGroupSnapshot,
			ID:              apimodels.ResourceID(resourceID),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.CloudWatchLogGroupSchema,
		})
	}

	return resources, nil
}
