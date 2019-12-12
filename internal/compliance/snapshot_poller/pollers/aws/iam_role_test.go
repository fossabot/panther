package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestIAMRolesList(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"ListRolesPages"})

	out := listRoles(mockSvc)
	assert.Equal(t, awstest.ExampleIAMRole, out[0])
}

func TestIAMRolesListError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"ListRolesPages"})

	out := listRoles(mockSvc)
	assert.Nil(t, out)
}

func TestIAMRolesGetPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GetRolePolicy"})

	out := getRolePolicy(mockSvc, aws.String("RoleName"), aws.String("PolicyName"))
	assert.NotEmpty(t, out)
}

func TestIAMRolesGetPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"ListRolesPages"})

	out := listRoles(mockSvc)
	assert.Nil(t, out)
}

func TestIAMRolesGetPolicies(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{
		"ListRolePoliciesPages",
		"ListAttachedRolePoliciesPages",
	})

	inlinePolicies, managedPolicies, err := getRolePolicies(mockSvc, aws.String("Franklin"))
	require.NoError(t, err)
	assert.Equal(
		t,
		[]*string{aws.String("AdministratorAccess")},
		managedPolicies,
	)
	assert.Equal(
		t,
		[]*string{aws.String("KinesisWriteOnly"), aws.String("SQSCreateQueue")},
		inlinePolicies,
	)
}

func TestIAMRolesGetPoliciesErrors(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{
		"ListRolePoliciesPages",
		"ListAttachedRolePoliciesPages",
	})

	inlinePolicies, managedPolicies, err := getRolePolicies(mockSvc, aws.String("Franklin"))
	require.Error(t, err)
	assert.Empty(t, inlinePolicies)
	assert.Empty(t, managedPolicies)
}

func TestIAMRolesPoller(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollIAMRoles(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
	assert.Len(t, resources, 1)
	assert.Equal(t, awstest.ExampleIAMRole.Arn, resources[0].Attributes.(*awsmodels.IAMRole).ARN)
}

func TestIAMRolesPollerError(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollIAMRoles(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Nil(t, resources)
}
