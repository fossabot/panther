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
