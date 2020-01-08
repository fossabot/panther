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

func TestS3GetBucketLogging(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketLogging"})

	out, err := getBucketLogging(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetBucketLoggingError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketLogging"})

	out, err := getBucketLogging(mockSvc, awstest.ExampleBucketName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestS3GetBucketTagging(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketTagging"})

	out, err := getBucketTagging(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetBucketTaggingError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketTagging"})

	out, err := getBucketTagging(mockSvc, awstest.ExampleBucketName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestS3GetBucketAcl(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketAcl"})

	out, err := getBucketACL(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetBucketAclError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketAcl"})

	out, err := getBucketACL(mockSvc, awstest.ExampleBucketName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestS3GetObjectLockConfiguration(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetObjectLockConfiguration"})

	out, err := getObjectLockConfiguration(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetObjectLockConfigurationError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetObjectLockConfiguration"})

	out, err := getObjectLockConfiguration(mockSvc, awstest.ExampleBucketName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestS3BucketsList(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"ListBuckets"})

	out := listBuckets(mockSvc)
	assert.NotEmpty(t, out)
}

func TestS3BucketsListError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"ListBuckets"})

	out := listBuckets(mockSvc)
	assert.Empty(t, out)
}

func TestS3GetBucketEncryption(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketEncryption"})

	out, err := getBucketEncryption(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetBucketEncryptionError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketEncryption"})

	out, err := getBucketEncryption(mockSvc, awstest.ExampleBucketName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestS3GetBucketPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketPolicy"})

	out, err := getBucketPolicy(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetBucketPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketPolicy"})

	out, err := getBucketPolicy(mockSvc, awstest.ExampleBucketName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestS3GetBucketVersioning(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketVersioning"})

	out, err := getBucketVersioning(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetBucketVersioningError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketVersioning"})

	out, err := getBucketVersioning(mockSvc, awstest.ExampleBucketName)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestS3GetBucketLocation(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketLocation"})

	out := getBucketLocation(mockSvc, awstest.ExampleBucketName)
	assert.Equal(t, "us-west-2", *out)
	assert.NotEmpty(t, out)
}

func TestS3GetBucketLocationError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketLocation"})

	out := getBucketLocation(mockSvc, awstest.ExampleBucketName)
	assert.Nil(t, out)
}

func TestS3GetBucketLifecycleConfiguration(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetBucketLifecycleConfiguration"})

	out, err := getBucketLifecycleConfiguration(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetPublicAccessBlock(t *testing.T) {
	mockSvc := awstest.BuildMockS3Svc([]string{"GetPublicAccessBlock"})

	out, err := getPublicAccessBlock(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestS3GetPublicAccessBlockOtherError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetPublicAccessBlock"})

	out, err := getPublicAccessBlock(mockSvc, awstest.ExampleBucketName)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestS3GetPublicAccessBlockDoesNotExist(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetPublicAccessBlockDoesNotExist"})

	out, err := getPublicAccessBlock(mockSvc, awstest.ExampleBucketName)
	require.NoError(t, err)
	assert.Nil(t, out)
}

func TestS3GetPublicAccessBlockAnotherAWSErr(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetPublicAccessBlockAnotherAWSErr"})

	out, err := getPublicAccessBlock(mockSvc, awstest.ExampleBucketName)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestS3GetBucketLifecycleConfigurationError(t *testing.T) {
	mockSvc := awstest.BuildMockS3SvcError([]string{"GetBucketLifecycleConfiguration"})

	out, err := getBucketLifecycleConfiguration(mockSvc, awstest.ExampleBucketName)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestS3BucketPoller(t *testing.T) {
	awstest.MockS3ForSetup = awstest.BuildMockS3SvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	S3ClientFunc = awstest.SetupMockS3

	resources, err := PollS3Buckets(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Equal(t, "arn:aws:s3:::unit-test-cloudtrail-bucket", string(resources[0].ID))
	assert.NotEmpty(t, resources)
}

func TestS3BucketPollerError(t *testing.T) {
	awstest.MockS3ForSetup = awstest.BuildMockS3SvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	S3ClientFunc = awstest.SetupMockS3

	resources, err := PollS3Buckets(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Empty(t, resources)
}
