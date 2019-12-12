package aws

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestEC2DescribeInstances(t *testing.T) {
	mockSvc := awstest.BuildMockEC2Svc([]string{"DescribeInstancesPages"})

	out := describeInstances(mockSvc)
	assert.NotEmpty(t, out)
}

func TestEC2DescribeInstancesError(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcError([]string{"DescribeInstancesPages"})

	out := describeInstances(mockSvc)
	assert.Nil(t, out)
}

func TestEC2BuildInstanceSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcAll()

	ec2Snapshot := buildEc2InstanceSnapshot(mockSvc, awstest.ExampleInstance)
	assert.NotEmpty(t, ec2Snapshot.SecurityGroups)
	assert.NotEmpty(t, ec2Snapshot.BlockDeviceMappings)
}

func TestEC2PollInstances(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2Instances(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Regexp(
		t,
		regexp.MustCompile(`arn:aws:ec2:.*:123456789012:instance/instance-aabbcc123`),
		resources[0].ID,
	)
	assert.NotEmpty(t, resources)
}

func TestEC2PollInstancesError(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2Instances(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Empty(t, resources)
}
