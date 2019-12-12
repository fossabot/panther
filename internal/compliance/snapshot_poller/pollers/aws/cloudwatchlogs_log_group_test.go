package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestCloudWatchLogsLogGroupsDescribe(t *testing.T) {
	mockSvc := awstest.BuildMockCloudWatchLogsSvc([]string{"DescribeLogGroupsPages"})

	out := describeLogGroups(mockSvc)
	assert.NotEmpty(t, out)
}

func TestCloudWatchLogsLogGroupsDescribeError(t *testing.T) {
	mockSvc := awstest.BuildMockCloudWatchLogsSvcError([]string{"DescribeLogGroupsPages"})

	out := describeLogGroups(mockSvc)
	assert.Nil(t, out)
}

func TestCloudWatchLogsLogGroupsListTags(t *testing.T) {
	mockSvc := awstest.BuildMockCloudWatchLogsSvc([]string{"ListTagsLogGroup"})

	out := listTagsLogGroup(mockSvc, awstest.ExampleDescribeLogGroups.LogGroups[0].LogGroupName)
	assert.NotEmpty(t, out)
}

func TestCloudWatchLogsLogGroupsListTagsError(t *testing.T) {
	mockSvc := awstest.BuildMockCloudWatchLogsSvcError([]string{"ListTagsLogGroup"})

	out := listTagsLogGroup(mockSvc, awstest.ExampleDescribeLogGroups.LogGroups[0].LogGroupName)
	assert.Nil(t, out)
}

func TestBuildCloudWatchLogsLogGroupSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockCloudWatchLogsSvcAll()

	certSnapshot := buildCloudWatchLogsLogGroupSnapshot(
		mockSvc,
		awstest.ExampleDescribeLogGroups.LogGroups[0],
	)

	assert.NotNil(t, certSnapshot.ARN)
	assert.NotNil(t, certSnapshot.StoredBytes)
	assert.Equal(t, "LogGroup-1", *certSnapshot.Name)
}

func TestCloudWatchLogsLogGroupPoller(t *testing.T) {
	awstest.MockCloudWatchLogsForSetup = awstest.BuildMockCloudWatchLogsSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	CloudWatchLogsClientFunc = awstest.SetupMockCloudWatchLogs

	resources, err := PollCloudWatchLogsLogGroups(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
}

func TestCloudWatchLogsLogGroupPollerError(t *testing.T) {
	awstest.MockCloudWatchLogsForSetup = awstest.BuildMockCloudWatchLogsSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	AcmClientFunc = awstest.SetupMockAcm

	resources, err := PollCloudWatchLogsLogGroups(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	for _, event := range resources {
		assert.Nil(t, event.Attributes)
	}
}
