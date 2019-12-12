package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

// Unit Tests

func TestEC2DescribeImages(t *testing.T) {
	mockSvc := awstest.BuildMockEC2Svc([]string{"DescribeImages"})
	ec2Amis = make(map[string][]*string)
	ec2Amis[defaultRegion] = []*string{awstest.ExampleAmi.ImageId}

	out, err := describeImages(mockSvc, defaultRegion)
	assert.NotEmpty(t, out)
	assert.NoError(t, err)
}

func TestEC2DescribeImagesError(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcError([]string{"DescribeImages"})

	out, err := describeImages(mockSvc, defaultRegion)
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestEC2BuildAmiSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcAll()

	ec2Snapshot := buildEc2AmiSnapshot(
		mockSvc,
		awstest.ExampleAmi,
	)

	assert.Equal(t, ec2Snapshot.ID, aws.String("ari-abc234"))
	assert.Equal(t, ec2Snapshot.ImageType, aws.String("ramdisk"))
}

func TestEC2PollAmis(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2Amis(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
}

func TestEC2PollAmiError(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2Amis(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Empty(t, resources)
}
