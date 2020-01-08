package remediation

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
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	policymodels "github.com/panther-labs/panther/api/gateway/analysis/models"
	processormodels "github.com/panther-labs/panther/api/gateway/remediation/models"
	"github.com/panther-labs/panther/api/gateway/resources/models"
	organizationmodels "github.com/panther-labs/panther/api/lambda/organization/models"
)

type mockLambdaClient struct {
	lambdaiface.LambdaAPI
	mock.Mock
}

func (m *mockLambdaClient) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*lambda.InvokeOutput), args.Error(1)
}

type mockRoundTripper struct {
	http.RoundTripper
	mock.Mock
}

func (m *mockRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	args := m.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}

var (
	input = &processormodels.RemediateResource{
		PolicyID:   "policyId",
		ResourceID: "resourceId",
	}

	remediation = &processormodels.Remediations{
		"AWS.S3.EnableBucketEncryption": map[string]interface{}{
			"SSEAlgorithm": "AES256",
		},
	}

	policy = &policymodels.Policy{
		AutoRemediationID: "AWS.S3.EnableBucketEncryption",
		AutoRemediationParameters: map[string]string{
			"SSEAlgorithm": "AES256",
		},
	}

	resourceAttributes = map[string]interface{}{
		"Region": "us-west-2",
	}

	resource = &models.Resource{
		Attributes: resourceAttributes,
	}
)

func init() {
	crossAccountRoleName = "crossAccountRoleName"
	sessionDurationSeconds = "60"
	organizationsAPI = "organizationsApi"
	resourcesServiceHostname = "resourcesServiceHostname"
	policiesServiceHostname = "policiesServiceHostname"
}

func TestRemediate(t *testing.T) {
	mockClient := &mockLambdaClient{}
	mockCreds := &credentials.Credentials{}

	mockRoundTripper := &mockRoundTripper{}
	httpClient = &http.Client{Transport: mockRoundTripper}

	mockRemediatorLambdaClient := &mockLambdaClient{}
	remediator := &Invoker{
		lambdaClient: mockRemediatorLambdaClient,
	}

	getCreds = func(c client.ConfigProvider, roleARN string, options ...func(*stscreds.AssumeRoleProvider)) *credentials.Credentials {
		return mockCreds
	}

	getLambda = func(p client.ConfigProvider, cfgs *aws.Config) lambdaiface.LambdaAPI {
		return mockClient
	}

	expectedPayload := Payload{
		RemediationID: string(policy.AutoRemediationID),
		Resource:      resourceAttributes,
		Parameters:    policy.AutoRemediationParameters,
	}
	expectedInput := LambdaInput{
		Action:  aws.String(remediationAction),
		Payload: expectedPayload,
	}
	expectedSerializedInput, _ := jsoniter.Marshal(expectedInput)

	expectedLambdaInput := &lambda.InvokeInput{
		FunctionName: aws.String("arn:aws:lambda:us-west-2:123456789012:function:function"),
		Payload:      expectedSerializedInput,
	}

	mockClient.On("Invoke", expectedLambdaInput).Return(&lambda.InvokeOutput{}, nil)

	expectedOrganizationsLambdaPayload := organizationmodels.LambdaInput{
		GetOrganization: &organizationmodels.GetOrganizationInput{},
	}
	expectedOrganizationsLambdaSerializedPayload, _ := jsoniter.Marshal(expectedOrganizationsLambdaPayload)

	expectedOrganizationsLambdaInput := &lambda.InvokeInput{
		FunctionName: aws.String("organizationsApi"),
		Payload:      expectedOrganizationsLambdaSerializedPayload,
	}

	organizationsOutput := &organizationmodels.GetOrganizationOutput{
		Organization: &organizationmodels.Organization{
			RemediationConfig: &organizationmodels.RemediationConfig{
				AwsRemediationLambdaArn: aws.String("arn:aws:lambda:us-west-2:123456789012:function:function"),
			},
		},
	}
	organizationsPayload, _ := jsoniter.Marshal(organizationsOutput)

	mockRemediatorLambdaClient.On("Invoke", expectedOrganizationsLambdaInput).Return(&lambda.InvokeOutput{Payload: organizationsPayload}, nil)
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(generateResponse(policy, http.StatusOK), nil).Once()
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(generateResponse(resource, http.StatusOK), nil).Once()

	// Clearing up the cache
	cache = nil

	result := remediator.Remediate(input)
	assert.NoError(t, result)

	// Verify cache has been populated
	expectedValue := &aws.Config{Region: aws.String("us-west-2"), Credentials: mockCreds}
	assert.Equal(t, expectedValue, cache)

	mockRemediatorLambdaClient.AssertExpectations(t)
	mockRoundTripper.AssertExpectations(t)
}

func TestRemediateLambdaError(t *testing.T) {
	mockClient := &mockLambdaClient{}
	mockCreds := &credentials.Credentials{}

	mockRoundTripper := &mockRoundTripper{}
	httpClient = &http.Client{Transport: mockRoundTripper}

	getCreds = func(c client.ConfigProvider, roleARN string, options ...func(*stscreds.AssumeRoleProvider)) *credentials.Credentials {
		return mockCreds
	}

	getLambda = func(p client.ConfigProvider, cfgs *aws.Config) lambdaiface.LambdaAPI {
		return mockClient
	}
	remediator := &Invoker{lambdaClient: mockClient}
	mockClient.On("Invoke", mock.Anything).Return(&lambda.InvokeOutput{}, errors.New("error"))
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(generateResponse(policy, http.StatusOK), nil).Once()
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(generateResponse(resource, http.StatusOK), nil).Once()

	result := remediator.Remediate(input)
	assert.Error(t, result)

	mockClient.AssertExpectations(t)
	mockRoundTripper.AssertExpectations(t)
}

func TestRemediateLambdaFunctionError(t *testing.T) {
	mockClient := &mockLambdaClient{}
	mockCreds := &credentials.Credentials{}

	mockRoundTripper := &mockRoundTripper{}
	httpClient = &http.Client{Transport: mockRoundTripper}

	getCreds = func(c client.ConfigProvider, roleARN string, options ...func(*stscreds.AssumeRoleProvider)) *credentials.Credentials {
		return mockCreds
	}

	getLambda = func(p client.ConfigProvider, cfgs *aws.Config) lambdaiface.LambdaAPI {
		return mockClient
	}

	lambdaOutput := &lambda.InvokeOutput{
		FunctionError: aws.String("LambdaError"),
	}

	mockClient.On("Invoke", mock.Anything).Return(lambdaOutput, nil)
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(generateResponse(policy, http.StatusOK), nil).Once()
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(generateResponse(resource, http.StatusOK), nil).Once()

	remediator := &Invoker{lambdaClient: mockClient}
	result := remediator.Remediate(input)
	assert.Error(t, result)

	mockClient.AssertExpectations(t)
	mockRoundTripper.AssertExpectations(t)
}

func TestGetRemediations(t *testing.T) {
	mockClient := &mockLambdaClient{}
	mockCreds := &credentials.Credentials{}

	mockRemediatorLambdaClient := &mockLambdaClient{}
	remediator := &Invoker{
		lambdaClient: mockRemediatorLambdaClient,
	}

	getCreds = func(c client.ConfigProvider, roleARN string, options ...func(*stscreds.AssumeRoleProvider)) *credentials.Credentials {
		return mockCreds
	}

	getLambda = func(p client.ConfigProvider, cfgs *aws.Config) lambdaiface.LambdaAPI {
		return mockClient
	}

	expectedInput := LambdaInput{Action: aws.String(listRemediationsAction)}
	expectedSerializedInput, _ := jsoniter.Marshal(expectedInput)

	expectedLambdaInput := &lambda.InvokeInput{
		FunctionName: aws.String("arn:aws:lambda:us-west-2:123456789012:function:function"),
		Payload:      expectedSerializedInput,
	}

	serializedRemediations := []byte("{\"AWS.S3.EnableBucketEncryption\": {\"SSEAlgorithm\": \"AES256\"}}")
	mockClient.On("Invoke", expectedLambdaInput).Return(&lambda.InvokeOutput{Payload: serializedRemediations}, nil)

	expectedOrganizationsLambdaPayload := organizationmodels.LambdaInput{
		GetOrganization: &organizationmodels.GetOrganizationInput{},
	}
	expectedOrganizationsLambdaSerializedPayload, _ := jsoniter.Marshal(expectedOrganizationsLambdaPayload)

	expectedOrganizationsLambdaInput := &lambda.InvokeInput{
		FunctionName: aws.String("organizationsApi"),
		Payload:      expectedOrganizationsLambdaSerializedPayload,
	}

	organizationsOutput := &organizationmodels.GetOrganizationOutput{
		Organization: &organizationmodels.Organization{
			RemediationConfig: &organizationmodels.RemediationConfig{
				AwsRemediationLambdaArn: aws.String("arn:aws:lambda:us-west-2:123456789012:function:function"),
			},
		},
	}
	organizationsPayload, _ := jsoniter.Marshal(organizationsOutput)

	mockRemediatorLambdaClient.On("Invoke", expectedOrganizationsLambdaInput).Return(&lambda.InvokeOutput{Payload: organizationsPayload}, nil)

	// Clearing up the cache
	cache = nil

	result, err := remediator.GetRemediations()
	assert.NoError(t, err)
	assert.Equal(t, remediation, result)
}

func TestDoNotTakeActionIfNoRemediationConfigured(t *testing.T) {
	mockRoundTripper := &mockRoundTripper{}
	httpClient = &http.Client{Transport: mockRoundTripper}

	mockRemediatorLambdaClient := &mockLambdaClient{}
	remediator := &Invoker{
		lambdaClient: mockRemediatorLambdaClient,
	}

	policy := &policymodels.Policy{
		AutoRemediationID: "",
	}

	mockRoundTripper.On("RoundTrip", mock.Anything).Return(generateResponse(policy, http.StatusOK), nil).Once()

	// Clearing up the cache
	cache = nil

	result := remediator.Remediate(input)
	assert.NoError(t, result)

	mockRoundTripper.AssertExpectations(t)
}

func generateResponse(body interface{}, httpCode int) *http.Response {
	serializedBody, _ := jsoniter.MarshalToString(body)
	return &http.Response{StatusCode: httpCode, Body: ioutil.NopCloser(strings.NewReader(serializedBody))}
}
