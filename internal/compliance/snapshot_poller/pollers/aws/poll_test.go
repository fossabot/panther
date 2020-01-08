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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

var (
	testGetCallerIdentityOutput = &sts.GetCallerIdentityOutput{
		Account: aws.String("111111111111"),
		Arn:     aws.String("arn:aws:iam::account-id:role/role-name"),
		UserId:  aws.String("mockUserId"),
	}
)

func configureMockSTSClientWithError(code, message string) *awstest.MockSTS {
	mockStsClient := &awstest.MockSTS{}
	mockStsClient.
		On("GetCallerIdentity", &sts.GetCallerIdentityInput{}).
		Return(
			testGetCallerIdentityOutput,
			awserr.New(code, message, errors.New("fake sts error")),
		)
	return mockStsClient
}

func configureMockSTSClientNoError() *awstest.MockSTS {
	mockStsClient := &awstest.MockSTS{}
	mockStsClient.
		On("GetCallerIdentity", &sts.GetCallerIdentityInput{}).
		Return(
			testGetCallerIdentityOutput,
			nil,
		)
	return mockStsClient
}

// Unit tests

func TestAssumeRole(t *testing.T) {
	AssumeRoleProviderFunc = awstest.STSAssumeRoleProviderMock
	STSClientFunc = awstest.SetupMockSTSClient
	awstest.MockSTSForSetup = configureMockSTSClientNoError()

	testSess, err := session.NewSession()
	require.NoError(t, err)

	creds, err := AssumeRole(
		&awsmodels.ResourcePollerInput{
			AuthSource:          &awstest.ExampleAuthSource,
			AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
			IntegrationID:       awstest.ExampleIntegrationID,
			Timestamp:           &awstest.ExampleTime,
		},
		testSess,
	)

	require.NoError(t, err)
	assert.NotEmpty(t, creds)
}

func TestAssumeRoleVerifyFailure(t *testing.T) {
	AssumeRoleProviderFunc = awstest.STSAssumeRoleProviderMock
	STSClientFunc = awstest.SetupMockSTSClient
	awstest.MockSTSForSetup = configureMockSTSClientWithError("AccessDenied", "You shall not pass")

	testSess, err := session.NewSession()
	require.NoError(t, err)

	creds, err := AssumeRole(
		&awsmodels.ResourcePollerInput{
			AuthSource:          &awstest.ExampleAuthSource,
			AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
			IntegrationID:       awstest.ExampleIntegrationID,
			Timestamp:           &awstest.ExampleTime,
		},
		testSess,
	)

	require.Error(t, err)
	assert.Nil(t, creds)
}

func TestAssumeRoleAddToCache(t *testing.T) {
	AssumeRoleProviderFunc = awstest.STSAssumeRoleProviderMock
	STSClientFunc = awstest.SetupMockSTSClient
	awstest.MockSTSForSetup = configureMockSTSClientNoError()

	// reset the cache for this test
	CredentialCache = make(map[string]*credentials.Credentials)

	testSess, err := session.NewSession()
	require.NoError(t, err)

	creds, err := AssumeRole(
		&awsmodels.ResourcePollerInput{
			AuthSource:          &awstest.ExampleAuthSource,
			AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
			IntegrationID:       aws.String("integration-id"),
			Timestamp:           &awstest.ExampleTime,
		},
		testSess,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, creds)

	creds, err = AssumeRole(
		&awsmodels.ResourcePollerInput{
			AuthSource:          &awstest.ExampleAuthSource2,
			AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
			IntegrationID:       aws.String("integration-id"),
			Timestamp:           &awstest.ExampleTime,
		},
		testSess,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, creds)

	assert.Len(t, CredentialCache, 2)
	assert.Contains(t, CredentialCache, awstest.ExampleAuthSource)
	assert.Contains(t, CredentialCache, awstest.ExampleAuthSource2)
}

func TestAssumeRoleNilSession(t *testing.T) {
	AssumeRoleProviderFunc = awstest.STSAssumeRoleProviderMock
	STSClientFunc = awstest.SetupMockSTSClient
	awstest.MockSTSForSetup = configureMockSTSClientNoError()

	creds, err := AssumeRole(
		&awsmodels.ResourcePollerInput{
			AuthSource:          &awstest.ExampleAuthSource,
			AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
			IntegrationID:       awstest.ExampleIntegrationID,
			Timestamp:           &awstest.ExampleTime,
		},
		nil,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, creds)
}

func TestAssumeRoleMissingParams(t *testing.T) {
	AssumeRoleProviderFunc = awstest.STSAssumeRoleProviderMock
	assert.Panics(t, func() { _, _ = AssumeRole(nil, nil) })
}

func TestVerifyAssumedCreds(t *testing.T) {
	STSClientFunc = awstest.SetupMockSTSClient
	awstest.MockSTSForSetup = configureMockSTSClientNoError()

	err := verifyAssumedCreds(&credentials.Credentials{})
	require.NoError(t, err)
}

func TestVerifyAssumedCredsAccessDeniedError(t *testing.T) {
	STSClientFunc = awstest.SetupMockSTSClient
	awstest.MockSTSForSetup = configureMockSTSClientWithError("AccessDenied", "You shall not pass")

	err := verifyAssumedCreds(&credentials.Credentials{})
	require.Error(t, err)
}

func TestVerifyAssumedCredsOtherError(t *testing.T) {
	STSClientFunc = awstest.SetupMockSTSClient
	awstest.MockSTSForSetup = configureMockSTSClientWithError("Error", "Something went wrong")

	err := verifyAssumedCreds(&credentials.Credentials{})
	// It's just logged
	require.NoError(t, err)
}
