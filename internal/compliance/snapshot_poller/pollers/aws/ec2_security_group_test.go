package aws

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestEC2DescribeSecurityGroups(t *testing.T) {
	mockSvc := awstest.BuildMockEC2Svc([]string{"DescribeSecurityGroupsPages"})

	out := describeSecurityGroups(mockSvc)
	assert.NotEmpty(t, out)
}

func TestEC2DescribeSecurityGroupsError(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcError([]string{"DescribeSecurityGroupsPages"})

	out := describeSecurityGroups(mockSvc)
	assert.Nil(t, out)
}

func TestEC2PollSecurityGroups(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2SecurityGroups(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Regexp(
		t,
		regexp.MustCompile(`arn:aws:ec2:.*:123456789012:security-group/sg-111222333`),
		resources[0].ID,
	)
	assert.NotEmpty(t, resources)
}

func TestEC2PollSecurityGroupsError(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2SecurityGroups(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Empty(t, resources)
}
