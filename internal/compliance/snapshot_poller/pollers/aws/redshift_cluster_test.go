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

func TestRedshiftClusterDescribe(t *testing.T) {
	mockSvc := awstest.BuildMockRedshiftSvc([]string{"DescribeClustersPages"})

	out := describeClusters(mockSvc)
	assert.NotEmpty(t, out)
}

func TestRedshiftClusterDescribeError(t *testing.T) {
	mockSvc := awstest.BuildMockRedshiftSvcError([]string{"DescribeClustersPages"})

	out := describeClusters(mockSvc)
	assert.Nil(t, out)
}

func TestRedshiftClusterDescribeLoggingStatus(t *testing.T) {
	mockSvc := awstest.BuildMockRedshiftSvc([]string{"DescribeLoggingStatus"})

	out, err := describeLoggingStatus(mockSvc, awstest.ExampleRDSSnapshotID)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRedshiftClusterDescribeLoggingStatusError(t *testing.T) {
	mockSvc := awstest.BuildMockRedshiftSvcError([]string{"DescribeLoggingStatus"})

	out, err := describeLoggingStatus(mockSvc, awstest.ExampleRDSSnapshotID)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestRedshiftClusterBuildSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockRedshiftSvcAll()

	clusterSnapshot := buildRedshiftClusterSnapshot(
		mockSvc,
		awstest.ExampleDescribeClustersOutput.Clusters[0],
	)

	assert.NotEmpty(t, clusterSnapshot.LoggingStatus)
	assert.NotEmpty(t, clusterSnapshot.GenericAWSResource)
}

func TestRedshiftCLusterBuildSnapshotErrors(t *testing.T) {
	mockSvc := awstest.BuildMockRedshiftSvcAllError()

	clusterSnapshot := buildRedshiftClusterSnapshot(
		mockSvc,
		awstest.ExampleDescribeClustersOutput.Clusters[0],
	)

	assert.NotEmpty(t, clusterSnapshot.GenericAWSResource)
	assert.NotNil(t, clusterSnapshot.VpcId)
	assert.Nil(t, clusterSnapshot.LoggingStatus)
}

func TestRedshiftClusterPoller(t *testing.T) {
	awstest.MockRedshiftForSetup = awstest.BuildMockRedshiftSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	RedshiftClientFunc = awstest.SetupMockRedshift

	resources, err := PollRedshiftClusters(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
	cluster := resources[0].Attributes.(*awsmodels.RedshiftCluster)
	assert.Equal(t, aws.String("awsuser"), cluster.MasterUsername)
	assert.Equal(t, aws.String("in-sync"), cluster.ClusterParameterGroups[0].ParameterApplyStatus)
	assert.Equal(t, aws.Int64(5439), cluster.Endpoint.Port)
	assert.Equal(t, aws.String("LEADER"), cluster.ClusterNodes[0].NodeRole)
	assert.False(t, *cluster.EnhancedVpcRouting)
}

func TestRedshiftClusterPollerError(t *testing.T) {
	awstest.MockRedshiftForSetup = awstest.BuildMockRedshiftSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	RedshiftClientFunc = awstest.SetupMockRedshift

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
