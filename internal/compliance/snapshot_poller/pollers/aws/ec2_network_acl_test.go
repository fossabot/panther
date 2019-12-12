package aws

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestEC2DescribeNetworkAcls(t *testing.T) {
	mockSvc := awstest.BuildMockEC2Svc([]string{"DescribeNetworkAclsPages"})

	out := describeNetworkAcls(mockSvc)
	assert.NotEmpty(t, out)
}

func TestEC2DescribeNetworkAclsError(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcError([]string{"DescribeNetworkAclsPages"})

	out := describeNetworkAcls(mockSvc)
	assert.Nil(t, out)
}

func TestEC2PollNetworkAcls(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2NetworkAcls(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Regexp(
		t,
		regexp.MustCompile(`arn:aws:ec2:.*:123456789012:network-acl/acl-111222333`),
		resources[0].ID,
	)
	assert.NotEmpty(t, resources)
}

func TestEC2PollNetworkAclsError(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2NetworkAcls(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Empty(t, resources)
}
