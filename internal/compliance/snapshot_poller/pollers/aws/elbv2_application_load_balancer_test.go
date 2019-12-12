package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestElbv2DescribeLoadBalancers(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2Svc([]string{"DescribeLoadBalancersPages"})

	out := describeLoadBalancers(mockSvc)
	assert.NotEmpty(t, out)
}

func TestElbv2DescribeLoadBalancersError(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2SvcError([]string{"DescribeLoadBalancersPages"})

	out := describeLoadBalancers(mockSvc)
	assert.Nil(t, out)
}

func TestElbv2DescribeListeners(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2Svc([]string{"DescribeListenersPages"})

	out := describeListeners(mockSvc, awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0].LoadBalancerArn)
	assert.NotEmpty(t, out)
}

func TestElbv2DescribeListenersError(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2SvcError([]string{"DescribeListenersPages"})

	out := describeListeners(mockSvc, awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0].LoadBalancerArn)
	assert.Nil(t, out)
}

func TestElbv2DescribeTags(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2Svc([]string{"DescribeTags"})

	out, err := describeTags(mockSvc, awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0].LoadBalancerArn)

	assert.Nil(t, err)
	assert.NotEmpty(t, out)
}

func TestElbv2DescribeTagsError(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2SvcError([]string{"DescribeTags"})

	out, err := describeTags(mockSvc, awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0].LoadBalancerArn)

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestElbv2DescribeSSLPolicies(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2Svc([]string{"DescribeSSLPolicies"})

	out, err := describeSSLPolicies(mockSvc)

	assert.Nil(t, err)
	assert.NotEmpty(t, out)
}

func TestElbv2DescribeSSLPoliciesError(t *testing.T) {
	mockSvc := awstest.BuildMockElbv2SvcError([]string{"DescribeSSLPolicies"})

	out, err := describeSSLPolicies(mockSvc)

	assert.Error(t, err)
	assert.Nil(t, out)
}
func TestBuildElbv2ApplicationLoadBalancerSnapshot(t *testing.T) {
	mockElbv2Svc := awstest.BuildMockElbv2SvcAll()
	mockWafRegionalSvc := awstest.BuildMockWafRegionalSvcAll()

	elbv2Snapshot := buildElbv2ApplicationLoadBalancerSnapshot(
		mockElbv2Svc,
		mockWafRegionalSvc,
		awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0],
	)

	assert.NotEmpty(t, elbv2Snapshot.SecurityGroups)
	assert.NotNil(t, elbv2Snapshot.WebAcl)
	assert.NotEmpty(t, elbv2Snapshot.Name)
}

func TestBuildElbv2ApplicationLoadBalancerSnapshotError(t *testing.T) {
	mockElbv2Svc := awstest.BuildMockElbv2SvcAllError()
	mockWafRegionalSvc := awstest.BuildMockWafRegionalSvcAllError()

	elbv2Snapshot := buildElbv2ApplicationLoadBalancerSnapshot(
		mockElbv2Svc,
		mockWafRegionalSvc,
		awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0],
	)

	assert.Nil(t, elbv2Snapshot.WebAcl)
	assert.Equal(
		t,
		awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0].LoadBalancerArn,
		elbv2Snapshot.ResourceID,
	)
}

func TestElbv2ApplicationLoadBalancersPoller(t *testing.T) {
	awstest.MockElbv2ForSetup = awstest.BuildMockElbv2SvcAll()
	awstest.MockWafRegionalForSetup = awstest.BuildMockWafRegionalSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	Elbv2ClientFunc = awstest.SetupMockElbv2
	WafRegionalClientFunc = awstest.SetupMockWafRegional

	resources, err := PollElbv2ApplicationLoadBalancers(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Equal(
		t,
		*awstest.ExampleDescribeLoadBalancersOutput.LoadBalancers[0].LoadBalancerArn,
		string(resources[0].ID),
	)
	assert.NotEmpty(t, resources[0].Attributes.(*awsmodels.Elbv2ApplicationLoadBalancer).Listeners)
	assert.NotNil(t, resources[0].Attributes.(*awsmodels.Elbv2ApplicationLoadBalancer).SSLPolicies)
	assert.NotNil(t, resources[0].Attributes.(*awsmodels.Elbv2ApplicationLoadBalancer).SSLPolicies["ELBSecurityPolicy1"])
	assert.NotEmpty(t, resources)
}

func TestElbv2ApplicationLoadBalancersPollerError(t *testing.T) {
	awstest.MockElbv2ForSetup = awstest.BuildMockElbv2SvcAllError()
	awstest.MockWafRegionalForSetup = awstest.BuildMockWafRegionalSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	Elbv2ClientFunc = awstest.SetupMockElbv2
	WafRegionalClientFunc = awstest.SetupMockWafRegional

	resources, err := PollElbv2ApplicationLoadBalancers(&awsmodels.ResourcePollerInput{
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
