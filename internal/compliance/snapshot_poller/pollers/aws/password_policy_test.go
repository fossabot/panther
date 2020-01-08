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

func TestGetPasswordPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GetAccountPasswordPolicy"})

	out, err := getPasswordPolicy(mockSvc)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestGetPasswordPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"GetAccountPasswordPolicy"})

	out, err := getPasswordPolicy(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestPasswordPolicyPoller(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvc([]string{"GetAccountPasswordPolicy"})

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollPasswordPolicy(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.Equal(t, "123456789012::AWS.PasswordPolicy", string(resources[0].ID))
}

func TestPasswordPolicyPollerError(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcError([]string{"GetAccountPasswordPolicy"})

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollPasswordPolicy(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Len(t, resources, 1)
}
