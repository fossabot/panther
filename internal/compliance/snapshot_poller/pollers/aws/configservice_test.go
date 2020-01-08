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
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestDescribeConfigurationRecorders(t *testing.T) {
	mockSvc := awstest.BuildMockConfigServiceSvc([]string{"DescribeConfigurationRecorders"})

	out, err := describeConfigurationRecorders(mockSvc)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestDescribeConfigurationRecordersError(t *testing.T) {
	mockSvc := awstest.BuildMockConfigServiceSvcError([]string{"DescribeConfigurationRecorders"})

	out, err := describeConfigurationRecorders(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestDescribeConfigurationRecorderStatus(t *testing.T) {
	mockSvc := awstest.BuildMockConfigServiceSvc([]string{"DescribeConfigurationRecorderStatus"})

	out, err := describeConfigurationRecorderStatus(mockSvc, awstest.ExampleConfigName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestDescribeConfigurationRecorderStatusError(t *testing.T) {
	mockSvc := awstest.BuildMockConfigServiceSvcError([]string{"DescribeConfigurationRecorderStatus"})

	out, err := describeConfigurationRecorderStatus(mockSvc, awstest.ExampleConfigName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestBuildConfigServiceSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockConfigServiceSvcAll()

	out := buildConfigServiceSnapshot(
		mockSvc,
		awstest.ExampleDescribeConfigurationRecorders.ConfigurationRecorders[0],
		"us-west-2",
	)
	assert.NotEmpty(t, out)
}

func TestBuildConfigServiceSnapshotError(t *testing.T) {
	mockSvc := awstest.BuildMockConfigServiceSvcAllError()

	out := buildConfigServiceSnapshot(
		mockSvc,
		awstest.ExampleDescribeConfigurationRecorders.ConfigurationRecorders[0],
		"us-west-2",
	)
	assert.NotEmpty(t, out)
}

func TestPollConfigServices(t *testing.T) {
	awstest.MockConfigServiceForSetup = awstest.BuildMockConfigServiceSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	ConfigServiceClientFunc = awstest.SetupMockConfigService

	resources, err := PollConfigServices(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)

	assert.IsType(t, &awsmodels.ConfigService{}, resources[0].Attributes)
	assert.Regexp(
		t, regexp.MustCompile(`123456789012\:.*\:AWS.Config.Recorder`), string(resources[0].ID),
	)

	assert.IsType(t, &awsmodels.ConfigServiceMeta{}, resources[len(resources)-1].Attributes)
	assert.Equal(t, "123456789012::AWS.Config.Recorder.Meta", string(resources[len(resources)-1].ID))
}

func TestPollConfigServicesError(t *testing.T) {
	awstest.MockConfigServiceForSetup = awstest.BuildMockConfigServiceSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	ConfigServiceClientFunc = awstest.SetupMockConfigService

	resources, err := PollConfigServices(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	// The meta resource should still send.
	assert.Len(t, resources, 1)
}
