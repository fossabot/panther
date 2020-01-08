package aws

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestIAMPolicyList(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"ListPoliciesPages"})

	out, err := listPolicies(mockSvc)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestIAMPolicyListError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"ListPoliciesPages"})

	out, err := listPolicies(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestIAMPolicyVersion(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GetPolicyVersion"})

	out, err := getPolicyVersion(
		mockSvc,
		aws.String("arn:aws:iam::aws:policy/aws-service-role/AWSSupportServiceRolePolicy"),
		aws.String("v2"),
	)

	require.NoError(t, err)
	assert.Equal(t, *awstest.ExamplePolicyDocumentDecoded, out)
	mockSvc.AssertExpectations(t)
}

func TestIAMPolicyVersionError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"GetPolicyVersion"})

	out, err := getPolicyVersion(
		mockSvc,
		aws.String("arn:aws:iam::aws:policy/aws-service-role/AWSSupportServiceRolePolicy"),
		aws.String("v2"),
	)

	require.NotNil(t, err)
	assert.Empty(t, out)
	mockSvc.AssertExpectations(t)
}

func TestIAMPolicyListEntities(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"ListEntitiesForPolicyPages"})

	out := listEntitiesForPolicy(
		mockSvc,
		aws.String("arn:aws:iam::aws:policy/aws-service-role/AWSSupportServiceRolePolicy"),
	)

	assert.NotEmpty(t, out)
}

func TestIAMPolicyListEntitiesError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"ListEntitiesForPolicyPages"})

	out := listEntitiesForPolicy(
		mockSvc,
		aws.String("arn:aws:iam::aws:policy/aws-service-role/AWSSupportServiceRolePolicy"),
	)

	assert.Empty(t, out)
}

func TestIAMPolicyBuildSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcAll()

	out := buildIAMPolicySnapshot(mockSvc, awstest.ExampleListPolicies.Policies[0])
	require.NotEmpty(t, out)
}

func TestIAMPolicyBuildSnapshotError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcAllError()

	out := buildIAMPolicySnapshot(mockSvc, awstest.ExampleListPolicies.Policies[0])
	require.Nil(t, out.Entities.PolicyGroups)
	require.Nil(t, out.Entities.PolicyRoles)
	require.Nil(t, out.Entities.PolicyUsers)
	assert.NotNil(t, out.PermissionsBoundaryUsageCount)
}

func TestIAMPolicyPoller(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollIamPolicies(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.Equal(t, *awstest.ExampleListPolicies.Policies[0].Arn, string(resources[0].ID))
}

func TestIAMPolicyPollerError(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollIamPolicies(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Nil(t, resources)
}
