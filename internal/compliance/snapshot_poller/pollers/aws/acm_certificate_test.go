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

func TestAcmCertificateList(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvc([]string{"ListCertificatesPages"})

	out := listCertificates(mockSvc)
	assert.NotEmpty(t, out)
}

func TestAcmCertificateListError(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvcError([]string{"ListCertificatesPages"})

	out := listCertificates(mockSvc)
	assert.Nil(t, out)
}

func TestAcmCertificateDescribe(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvc([]string{"DescribeCertificate"})

	out, err := describeCertificate(mockSvc, awstest.ExampleCertificateArn)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestAcmCertificateDescribeError(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvcError([]string{"DescribeCertificate"})

	out, err := describeCertificate(mockSvc, awstest.ExampleCertificateArn)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestAcmCertificateListTags(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvc([]string{"ListTagsForCertificate"})

	out, err := listTagsForCertificate(mockSvc, awstest.ExampleCertificateArn)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestAcmCertificateListTagsError(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvcError([]string{"ListTagsForCertificate"})

	out, err := listTagsForCertificate(mockSvc, awstest.ExampleCertificateArn)
	require.Error(t, err)
	assert.Nil(t, out)
}

func TestAcmCertificateBuildSnapshot(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvcAll()

	certSnapshot := buildAcmCertificateSnapshot(
		mockSvc,
		awstest.ExampleListCertificatesOutput.CertificateSummaryList[0].CertificateArn,
	)

	assert.NotEmpty(t, certSnapshot.ARN)
	assert.Equal(t, "Value1", *certSnapshot.Tags["Key1"])
}

func TestAcmCertificateBuildSnapshotErrors(t *testing.T) {
	mockSvc := awstest.BuildMockAcmSvcAllError()

	certSnapshot := buildAcmCertificateSnapshot(
		mockSvc,
		awstest.ExampleListCertificatesOutput.CertificateSummaryList[0].CertificateArn,
	)

	assert.Nil(t, certSnapshot)
}

func TestAcmCertificatePoller(t *testing.T) {
	awstest.MockAcmForSetup = awstest.BuildMockAcmSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	AcmClientFunc = awstest.SetupMockAcm

	resources, err := PollAcmCertificates(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Regions:             awstest.ExampleRegions,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Equal(t, *awstest.ExampleCertificateArn, string(resources[0].ID))
	assert.NotEmpty(t, resources)
}

func TestAcmCertificatePollerError(t *testing.T) {
	awstest.MockAcmForSetup = awstest.BuildMockAcmSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	AcmClientFunc = awstest.SetupMockAcm

	resources, err := PollAcmCertificates(&awsmodels.ResourcePollerInput{
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
