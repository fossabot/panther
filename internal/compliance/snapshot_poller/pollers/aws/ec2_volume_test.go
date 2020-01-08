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

func TestEC2DescribeVolumes(t *testing.T) {
	mockSvc := awstest.BuildMockEC2Svc([]string{"DescribeVolumesPages"})

	out := describeVolumes(mockSvc)
	assert.NotEmpty(t, out)
}

func TestEC2DescribeVolumesError(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcError([]string{"DescribeVolumesPages"})

	out := describeVolumes(mockSvc)
	assert.Nil(t, out)
}

func TestEC2DescribeSnapshots(t *testing.T) {
	mockSvc := awstest.BuildMockEC2Svc([]string{"DescribeSnapshotsPages"})

	out := describeSnapshots(mockSvc, awstest.ExampleVolumeId)
	assert.NotEmpty(t, out)
	assert.Len(t, out, 1)
}

func TestEC2DescribeSnapshotsError(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcError([]string{"DescribeSnapshotsPages"})

	out := describeSnapshots(mockSvc, awstest.ExampleVolumeId)
	assert.Nil(t, out)
}

func TestEC2DescribeSnapshotAttribute(t *testing.T) {
	mockSvc := awstest.BuildMockEC2Svc([]string{"DescribeSnapshotAttribute"})

	out, err := describeSnapshotAttribute(mockSvc, awstest.ExampleSnapshotId)
	assert.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Len(t, out, 1)
}

func TestEC2DescribeSnapshotAttributeError(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcError([]string{"DescribeSnapshotAttribute"})

	out, err := describeSnapshotAttribute(mockSvc, awstest.ExampleSnapshotId)
	assert.Nil(t, out)
	assert.Error(t, err)
}

func TestBuildEc2VolumeSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockEC2SvcAll()

	volumeSnapshot := buildEc2VolumeSnapshot(
		mockSvc,
		awstest.ExampleDescribeVolumesOutput.Volumes[0],
	)

	assert.NotNil(t, volumeSnapshot.AvailabilityZone)
	assert.NotEmpty(t, volumeSnapshot.Attachments)
}

func TestEc2VolumePoller(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2Volumes(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
}

func TestEc2VolumePollerError(t *testing.T) {
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	EC2ClientFunc = awstest.SetupMockEC2

	resources, err := PollEc2Volumes(&awsmodels.ResourcePollerInput{
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
