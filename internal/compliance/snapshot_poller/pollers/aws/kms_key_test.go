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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestKMSKeyList(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvc([]string{"ListKeys"})

	out := listKeys(mockSvc)
	assert.NotEmpty(t, out)
}

func TestKMSKeyListError(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvcError([]string{"ListKeys"})

	out := listKeys(mockSvc)
	assert.Nil(t, out)
}

func TestKMSKeyGetRotationStatus(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvc([]string{"GetKeyRotationStatus"})

	out, err := getKeyRotationStatus(mockSvc, awstest.ExampleKeyId)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestKMSKeyGetRotationStatusError(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvcError([]string{"GetKeyRotationStatus"})

	out, err := getKeyRotationStatus(mockSvc, awstest.ExampleKeyId)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestKMSKeyDescribe(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvc([]string{"DescribeKey"})

	out, err := describeKey(mockSvc, awstest.ExampleKeyId)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestKMSKeyDescribeError(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvcError([]string{"DescribeKey"})

	out, err := describeKey(mockSvc, awstest.ExampleKeyId)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestKMSKeyGetPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvc([]string{"GetKeyPolicy"})

	out, err := getKeyPolicy(mockSvc, awstest.ExampleKeyId)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestKMSKeyGetPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvcError([]string{"GetKeyPolicy"})

	out, err := getKeyPolicy(mockSvc, awstest.ExampleKeyId)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestKMSKeyListResourceTags(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvc([]string{"ListResourceTags"})

	out, err := listResourceTags(mockSvc, awstest.ExampleKeyId)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestKMSKeyListResourceTagsError(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvcError([]string{"ListResourceTags"})

	out, err := listResourceTags(mockSvc, awstest.ExampleKeyId)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestBuildKmsKeySnapshotAWSManaged(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvc([]string{
		"GetKeyRotationStatus",
		"GetKeyPolicy",
		"ListResourceTags",
	})
	// Return the AWS managed example
	mockSvc.
		On("DescribeKey", mock.Anything).
		Return(awstest.ExampleDescribeKeyOutputAWSManaged, nil)
	awstest.MockKmsForSetup = mockSvc

	keySnapshot := buildKmsKeySnapshot(mockSvc, awstest.ExampleListKeysOutput.Keys[0])
	assert.Nil(t, keySnapshot.KeyRotationEnabled)
	assert.NotEmpty(t, keySnapshot.KeyManager)
	assert.NotEmpty(t, keySnapshot.Policy)
}

func TestBuildKmsKeySnapshotErrors(t *testing.T) {
	mockSvc := awstest.BuildMockKmsSvcAllError()

	keySnapshot := buildKmsKeySnapshot(mockSvc, awstest.ExampleListKeysOutput.Keys[0])
	assert.Nil(t, keySnapshot)
}

func TestKMSKeyPoller(t *testing.T) {
	awstest.MockKmsForSetup = awstest.BuildMockKmsSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	KmsClientFunc = awstest.SetupMockKms

	resources, err := PollKmsKeys(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
}

func TestKMSKeyPollerError(t *testing.T) {
	awstest.MockKmsForSetup = awstest.BuildMockKmsSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	KmsClientFunc = awstest.SetupMockKms

	resources, err := PollKmsKeys(&awsmodels.ResourcePollerInput{
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
