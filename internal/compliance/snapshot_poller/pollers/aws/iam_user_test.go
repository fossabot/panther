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
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

func TestGetCredentialReport(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GetCredentialReport"})

	out, err := getCredentialReport(mockSvc)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestGetCredentialReportError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"GetCredentialReport"})

	out, err := getCredentialReport(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestGenerateCredentialReport(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GenerateCredentialReport"})

	out, err := generateCredentialReport(mockSvc)
	require.NoError(t, err)
	assert.Equal(t, awstest.ExampleGenerateCredentialReport, out)
}

func TestGenerateCredentialReportError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"GenerateCredentialReport"})

	out, err := generateCredentialReport(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestBuildCredentialReport(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{
		"GenerateCredentialReport",
		"GetCredentialReport",
	})

	credentialReport, err := buildCredentialReport(mockSvc)

	require.NoError(t, err)
	assert.Equal(t, awstest.ExampleExtractedCredentialReport, credentialReport)
}

func TestBuildCredentialReportError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{
		"GenerateCredentialReport",
		"GetCredentialReport",
	})

	credentialReport, err := buildCredentialReport(mockSvc)

	require.NotNil(t, err)
	assert.Empty(t, credentialReport)
}

func TestGenerateCredentialReportInProgress(t *testing.T) {
	awstest.GenerateCredentialReportInProgress = true
	defer func() { awstest.GenerateCredentialReportInProgress = false }()
	mockSvc := awstest.BuildMockIAMSvc([]string{"GenerateCredentialReport"})

	out, err := generateCredentialReport(mockSvc)
	require.NoError(t, err)
	assert.Equal(t, awstest.ExampleGenerateCredentialReport, out)
}

func TestExtractCredentialReport(t *testing.T) {
	out, err := extractCredentialReport(awstest.ExampleCredentialReport.Content)
	require.NoError(t, err)
	assert.Equal(t, awstest.ExampleExtractedCredentialReport, out)
}

func TestExtractCredentialReportError(t *testing.T) {
	badReport := []byte("h1,h2,h3\nv1,v2\nv1,v2,v3,v4,v5\n")
	out, err := extractCredentialReport(badReport)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestIAMUsersList(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"ListUsersPages"})

	out := listUsers(mockSvc)
	assert.Equal(t, awstest.ExampleListUsers.Users, out)
}

func TestIAMUsersListError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"ListUsersPages"})

	out := listUsers(mockSvc)
	assert.Nil(t, out)
}

func TestIAMUsersGetPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GetUserPolicy"})

	out := getUserPolicy(mockSvc, aws.String("ExampleUser"), aws.String("ExamplePolicy"))
	assert.Equal(t, awstest.ExampleGetUserPolicy.PolicyDocument, out)
}

func TestIAMUsersGetUserPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"GetUserPolicy"})

	out := getUserPolicy(mockSvc, aws.String("ExampleUser"), aws.String("ExamplePolicy"))
	assert.Nil(t, out)
}

func TestIAMUsersGetPolicies(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{
		"ListUserPoliciesPages",
		"ListAttachedUserPoliciesPages",
	})

	inlinePolicies, managedPolicies, err := getUserPolicies(mockSvc, aws.String("Franklin"))
	require.NoError(t, err)
	assert.Equal(
		t,
		[]*string{aws.String("ForceMFA"), aws.String("IAMAdministrator")},
		managedPolicies,
	)
	assert.Equal(
		t,
		[]*string{aws.String("KinesisWriteOnly"), aws.String("SQSCreateQueue")},
		inlinePolicies,
	)
}

func TestIAMUsersGetPoliciesErrors(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{
		"ListUserPoliciesPages",
		"ListAttachedUserPoliciesPages",
	})

	inlinePolicies, managedPolicies, err := getUserPolicies(mockSvc, aws.String("Franklin"))
	require.Error(t, err)
	assert.Empty(t, inlinePolicies)
	assert.Empty(t, managedPolicies)
}

func TestIAMUsersListVirtualMFADevices(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"ListVirtualMFADevicesPages"})

	expected := map[string]*awsmodels.VirtualMFADevice{
		"123456789012": {
			EnableDate:   awstest.ExampleDate,
			SerialNumber: aws.String("arn:aws:iam::123456789012:mfa/root-account-mfa-device"),
		},
		"AAAAAAAQQQQQO2HVVVVVV": {
			EnableDate:   awstest.ExampleDate,
			SerialNumber: aws.String("arn:aws:iam::123456789012:mfa/unit_test_user"),
		},
	}

	out, err := listVirtualMFADevices(mockSvc)
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestIAMUsersListVirtualMFADevicesError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"ListVirtualMFADevicesPages"})

	out, err := listVirtualMFADevices(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestIAMUsersPoller(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcAll()

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollIAMUsers(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	rootSnapshot := &awsmodels.IAMRootUser{
		GenericResource: awsmodels.GenericResource{
			ResourceID:   aws.String("arn:aws:iam::123456789012:root"),
			ResourceType: aws.String(awsmodels.IAMRootUserSchema),
			TimeCreated:  utils.DateTimeFormat(*awstest.ExampleDate),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			AccountID: awstest.ExampleAccountId,
			ARN:       aws.String("arn:aws:iam::123456789012:root"),
			ID:        aws.String("123456789012"),
			Name:      aws.String("<root_account>"),
			Region:    aws.String(awsmodels.GlobalRegion),
		},
		CredentialReport: awstest.ExampleExtractedCredentialReport["<root_account>"],
		VirtualMFA: &awsmodels.VirtualMFADevice{
			EnableDate:   awstest.ExampleDate,
			SerialNumber: aws.String("arn:aws:iam::123456789012:mfa/root-account-mfa-device"),
		},
	}

	expectedIamUserSnapshots := []*awsmodels.IAMUser{
		{
			GenericResource: awsmodels.GenericResource{
				ResourceID:   aws.String("arn:aws:iam::123456789012:user/unit_test_user"),
				ResourceType: aws.String(awsmodels.IAMUserSchema),
				TimeCreated:  utils.DateTimeFormat(*awstest.ExampleDate),
			},
			GenericAWSResource: awsmodels.GenericAWSResource{
				AccountID: awstest.ExampleAccountId,
				ARN:       aws.String("arn:aws:iam::123456789012:user/unit_test_user"),
				ID:        aws.String("AAAAAAAQQQQQO2HVVVVVV"),
				Name:      aws.String("unit_test_user"),
				Region:    aws.String(awsmodels.GlobalRegion),
				Tags:      map[string]*string{},
			},
			Path:             aws.String("/service_accounts/"),
			CredentialReport: awstest.ExampleExtractedCredentialReport["unit_test_user"],
			Groups: []*iam.Group{
				awstest.ExampleGroup,
			},
			VirtualMFA: &awsmodels.VirtualMFADevice{
				EnableDate:   awstest.ExampleDate,
				SerialNumber: aws.String("arn:aws:iam::123456789012:mfa/unit_test_user"),
			},
			InlinePolicies: map[string]*string{
				"KinesisWriteOnly": aws.String("JSON POLICY DOCUMENT"),
				"SQSCreateQueue":   aws.String("JSON POLICY DOCUMENT"),
			},
			ManagedPolicyNames: []*string{aws.String("ForceMFA"), aws.String("IAMAdministrator")},
		},
		{
			GenericResource: awsmodels.GenericResource{
				ResourceID:   aws.String("arn:aws:iam::123456789012:user/Franklin"),
				ResourceType: aws.String(awsmodels.IAMUserSchema),
				TimeCreated:  utils.DateTimeFormat(*awstest.ExampleDate),
			},
			GenericAWSResource: awsmodels.GenericAWSResource{
				AccountID: awstest.ExampleAccountId,
				ARN:       aws.String("arn:aws:iam::123456789012:user/Franklin"),
				ID:        aws.String("AIDA4PIQ2YYOO2HYP2JNV"),
				Name:      aws.String("Franklin"),
				Region:    aws.String(awsmodels.GlobalRegion),
				Tags:      map[string]*string{},
			},
			Path:             aws.String("/"),
			CredentialReport: awstest.ExampleExtractedCredentialReport["Franklin"],
			Groups: []*iam.Group{
				awstest.ExampleGroup,
			},
			VirtualMFA: nil,
			InlinePolicies: map[string]*string{
				"KinesisWriteOnly": aws.String("JSON POLICY DOCUMENT"),
				"SQSCreateQueue":   aws.String("JSON POLICY DOCUMENT"),
			},
			ManagedPolicyNames: []*string{aws.String("ForceMFA"), aws.String("IAMAdministrator")},
		},
	}

	require.NoError(t, err)
	assert.NotEmpty(t, resources)
	// Root and two IAM users
	assert.Len(t, resources, 3)
	assert.Equal(t, expectedIamUserSnapshots[0], resources[0].Attributes)
	assert.Equal(t, expectedIamUserSnapshots[1], resources[1].Attributes)
	assert.Equal(t, rootSnapshot, resources[2].Attributes)
}

func TestIAMUsersPollerError(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcAllError()

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollIAMUsers(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	// Even though ListUsers will return no users, the poller will continue making API calls to build
	// the root user and error when building the credential report
	require.Error(t, err)
	assert.Nil(t, resources)
}
