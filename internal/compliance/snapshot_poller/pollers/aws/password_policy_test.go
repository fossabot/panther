package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestGetPasswordPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GetAccountPasswordPolicy"})

	out, err := getPasswordPolicy(mockSvc)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestGetPasswordPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"GetAccountPasswordPolicy"})

	out, err := getPasswordPolicy(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestPasswordPolicyPoller(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvc([]string{"GetAccountPasswordPolicy"})

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollPasswordPolicy(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.Equal(t, "123456789012::AWS.PasswordPolicy", string(resources[0].ID))
}

func TestPasswordPolicyPollerError(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcError([]string{"GetAccountPasswordPolicy"})

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollPasswordPolicy(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Len(t, resources, 1)
}
