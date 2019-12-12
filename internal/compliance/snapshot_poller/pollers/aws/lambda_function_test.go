package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestLambdaFunctionsList(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvc([]string{"ListFunctionsPages"})

	out := listFunctions(mockSvc)
	assert.NotEmpty(t, out)
}

func TestLambdaFunctionsListError(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvcError([]string{"ListFunctionsPages"})

	out := listFunctions(mockSvc)
	assert.Nil(t, out)
}

func TestLambdaFunctionListTags(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvc([]string{"ListTags"})

	out, err := listTagsLambda(mockSvc, awstest.ExampleFunctionName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestLambdaFunctionListTagsError(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvcError([]string{"ListTags"})

	out, err := listTagsLambda(mockSvc, awstest.ExampleFunctionName)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestLambdaFunctionGetPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvc([]string{"GetPolicy"})

	out, err := getPolicy(mockSvc, awstest.ExampleFunctionName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestLambdaFunctionGetPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvcError([]string{"GetPolicy"})

	out, err := getPolicy(mockSvc, awstest.ExampleFunctionName)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestBuildLambdaFunctionSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvcAll()

	lambdaSnapshot := buildLambdaFunctionSnapshot(
		mockSvc,
		awstest.ExampleListFunctions.Functions[0],
	)

	assert.NotEmpty(t, lambdaSnapshot.Tags)
	assert.NotEmpty(t, lambdaSnapshot.Policy)
	assert.Equal(t, "arn:aws:lambda:us-west-2:123456789012:function:ExampleFunction", *lambdaSnapshot.ARN)
	assert.Equal(t, awstest.ExampleFunctionConfiguration.TracingConfig, lambdaSnapshot.TracingConfig)
}

func TestBuildLambdaFunctionSnapshotErrors(t *testing.T) {
	mockSvc := awstest.BuildMockLambdaSvcAllError()

	lambdaSnapshot := buildLambdaFunctionSnapshot(
		mockSvc,
		awstest.ExampleListFunctions.Functions[0],
	)

	assert.NotNil(t, lambdaSnapshot)
	assert.Nil(t, lambdaSnapshot.Policy)
	assert.Nil(t, lambdaSnapshot.Tags)
}

func TestLambdaFunctionPoller(t *testing.T) {
	awstest.MockLambdaForSetup = awstest.BuildMockLambdaSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	LambdaClientFunc = awstest.SetupMockLambda

	resources, err := PollLambdaFunctions(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
}

func TestLambdaFunctionPollerError(t *testing.T) {
	awstest.MockLambdaForSetup = awstest.BuildMockLambdaSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	LambdaClientFunc = awstest.SetupMockLambda

	resources, err := PollLambdaFunctions(&awsmodels.ResourcePollerInput{
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
