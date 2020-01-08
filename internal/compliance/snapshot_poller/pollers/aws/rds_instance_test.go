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

func TestRDSInstanceDescribe(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvc([]string{"DescribeDBInstancesPages"})

	out := describeDBInstances(mockSvc)
	assert.NotEmpty(t, out)
}

func TestRDSInstanceDescribeError(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvcError([]string{"DescribeDBInstancesPages"})

	out := describeDBInstances(mockSvc)
	assert.Nil(t, out)
}

func TestRDSInstanceDescribeSnapshots(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvc([]string{"DescribeDBSnapshotsPages"})

	out, err := describeDBSnapshots(mockSvc, awstest.ExampleRDSInstanceName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRDSInstanceDescribeSnapshotsError(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvcError([]string{"DescribeDBSnapshotsPages"})

	out, err := describeDBSnapshots(mockSvc, awstest.ExampleRDSInstanceName)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestRDSInstanceListTagsForResource(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvc([]string{"ListTagsForResource"})

	out, err := listTagsForResourceRds(mockSvc, awstest.ExampleRDSInstanceName)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRDSInstanceListTagsForResourceError(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvcError([]string{"ListTagsForResource"})

	out, err := listTagsForResourceRds(mockSvc, awstest.ExampleRDSInstanceName)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestRDSInstanceDescribeSnapshotAttributes(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvc([]string{"DescribeDBSnapshotAttributes"})

	out, err := describeDBSnapshotAttributes(mockSvc, awstest.ExampleRDSSnapshotID)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRDSInstanceDescribeSnapshotAttributesError(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvcError([]string{"DescribeDBSnapshotAttributes"})

	out, err := describeDBSnapshotAttributes(mockSvc, awstest.ExampleRDSSnapshotID)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestRDSInstanceBuildSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvcAll()

	instanceSnapshot := buildRDSInstanceSnapshot(
		mockSvc,
		awstest.ExampleDescribeDBInstancesOutput.DBInstances[0],
	)

	assert.NotEmpty(t, instanceSnapshot.ARN)
	assert.NotEmpty(t, instanceSnapshot.SnapshotAttributes)
}

func TestRDSInstanceBuildSnapshotErrors(t *testing.T) {
	mockSvc := awstest.BuildMockRdsSvcAllError()

	instance := buildRDSInstanceSnapshot(
		mockSvc,
		awstest.ExampleDescribeDBInstancesOutput.DBInstances[0],
	)

	assert.Equal(t, "db.t2.micro", *instance.DBInstanceClass)
	assert.Equal(t, *awstest.ExampleRDSInstanceName, *instance.ID)
	assert.Equal(t, awstest.ExampleDescribeDBInstancesOutput.DBInstances[0].OptionGroupMemberships, instance.OptionGroupMemberships)
}

func TestRDSInstancePoller(t *testing.T) {
	awstest.MockRdsForSetup = awstest.BuildMockRdsSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	RDSClientFunc = awstest.SetupMockRds

	resources, err := PollRDSInstances(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
	instance := resources[0].Attributes.(*awsmodels.RDSInstance)
	assert.Equal(t, aws.String("superuser"), instance.MasterUsername)
	assert.Equal(t,
		aws.String("restore"),
		instance.SnapshotAttributes[0].DBSnapshotAttributes[0].AttributeName,
	)
	assert.Equal(t, aws.Int64(3306), instance.Endpoint.Port)
	assert.NotEmpty(t, instance.DBSubnetGroup.Subnets)
	assert.Equal(t, aws.String("in-sync"), instance.OptionGroupMemberships[0].Status)
}

func TestRDSInstancePollerError(t *testing.T) {
	awstest.MockRdsForSetup = awstest.BuildMockRdsSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	RDSClientFunc = awstest.SetupMockRds

	resources, err := PollRDSInstances(&awsmodels.ResourcePollerInput{
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
