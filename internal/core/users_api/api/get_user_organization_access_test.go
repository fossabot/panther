package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	organizationmodels "github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/api/lambda/users/models"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
	"github.com/panther-labs/panther/pkg/genericapi"
)

func TestGetUserOrganizationAccessSuccess(t *testing.T) {
	m := &users.MockTable{}
	organizationAPI = "organizationAPI"
	userTable = m
	email := aws.String("joe@blow.com")
	m.On("Get", email).Return(&models.UserItem{
		ID: aws.String("user-123"),
	}, nil)

	ml := &mockLambdaClient{}
	lambdaClient = ml
	getOrgOutput := organizationmodels.GetOrganizationOutput{
		Organization: &organizationmodels.Organization{
			CreatedAt:            aws.String("time-123"),
			DisplayName:          aws.String("Initech"),
			Email:                aws.String("joe@blow.com"),
			AlertReportFrequency: aws.String("errday"),
			AwsConfig: &organizationmodels.AwsConfig{
				UserPoolID:     aws.String("userpool-123"),
				AppClientID:    aws.String("client-123"),
				IdentityPoolID: aws.String("identitypool-123"),
			},
		},
	}
	mockOrgLambdaResponsePayload, err := jsoniter.Marshal(getOrgOutput)
	require.NoError(t, err)
	mockOrgLambdaResponse := &lambda.InvokeOutput{Payload: mockOrgLambdaResponsePayload}
	expecteOrgLambdaPayload, err := jsoniter.Marshal(
		organizationmodels.LambdaInput{GetOrganization: &organizationmodels.GetOrganizationInput{}})
	require.NoError(t, err)
	expectedOrgLambdaInput := &lambda.InvokeInput{FunctionName: aws.String("organizationAPI"), Payload: expecteOrgLambdaPayload}
	ml.On("Invoke", expectedOrgLambdaInput).Return(mockOrgLambdaResponse, nil)

	result, err := (API{}).GetUserOrganizationAccess(&models.GetUserOrganizationAccessInput{
		Email: email,
	})
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestGetUserOrganizationAccessGetUserFailed(t *testing.T) {
	m := &users.MockTable{}
	organizationAPI = "organizationAPI"
	userTable = m
	email := aws.String("joe@blow.com")
	m.On("Get", email).Return(&models.UserItem{}, &genericapi.LambdaError{})

	result, err := (API{}).GetUserOrganizationAccess(&models.GetUserOrganizationAccessInput{
		Email: email,
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetUserOrganizationAccessGetOrganizationsFailed(t *testing.T) {
	m := &users.MockTable{}
	organizationAPI = "organizationAPI"
	userTable = m
	email := aws.String("joe@blow.com")
	m.On("Get", email).Return(&models.UserItem{
		ID: aws.String("user-123"),
	}, nil)

	ml := &mockLambdaClient{}
	lambdaClient = ml
	getOrgOutput := organizationmodels.GetOrganizationOutput{}
	mockOrgLambdaResponsePayload, err := jsoniter.Marshal(getOrgOutput)
	require.NoError(t, err)
	mockOrgLambdaResponse := &lambda.InvokeOutput{Payload: mockOrgLambdaResponsePayload}
	expectedOrgLambdaPayload, err := jsoniter.Marshal(
		organizationmodels.LambdaInput{GetOrganization: &organizationmodels.GetOrganizationInput{}})
	require.NoError(t, err)
	expectedOrgLambdaInput := &lambda.InvokeInput{FunctionName: aws.String("organizationAPI"), Payload: expectedOrgLambdaPayload}
	ml.On("Invoke", expectedOrgLambdaInput).Return(mockOrgLambdaResponse, &genericapi.LambdaError{})

	result, err := (API{}).GetUserOrganizationAccess(&models.GetUserOrganizationAccessInput{
		Email: email,
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}
