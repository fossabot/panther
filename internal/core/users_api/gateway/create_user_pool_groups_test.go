package gateway

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/pkg/genericapi"
)

func TestCreateGroups(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)
	mockCognitoClient.On(
		"CreateGroup", mock.Anything).Return(&provider.CreateGroupOutput{}, nil)

	assert.NoError(t, gw.CreateUserPoolGroups(aws.String("us-west-2_ZlG7Ldp1K")))
	mockCognitoClient.AssertExpectations(t)
	mockCognitoClient.AssertNumberOfCalls(t, "CreateGroup", 3)
}

func TestCreateGroupsFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockIamService.On(
		"GetRole", mock.Anything).Return(&iam.GetRoleOutput{
		Role: &iam.Role{
			Arn: aws.String("fakeArn"),
		},
	}, nil)
	mockCognitoClient.On(
		"CreateGroup", mock.Anything).Return(&provider.CreateGroupOutput{}, &genericapi.AWSError{})

	assert.Error(t, gw.CreateUserPoolGroups(aws.String("us-west-2_ZlG7Ldp1K")))
	mockCognitoClient.AssertExpectations(t)
	mockCognitoClient.AssertNumberOfCalls(t, "CreateGroup", 1)
}

func TestGetAdminRoleFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}
	mockRoleParam := iam.GetRoleInput{
		RoleName: aws.String(IdentityPoolAuthenticatedAdminsRole),
	}
	mockIamService.On(
		"GetRole", &mockRoleParam).Return(&iam.GetRoleOutput{}, &genericapi.AWSError{})

	assert.Error(t, gw.CreateUserPoolGroups(aws.String("us-west-2_ZlG7Ldp1K")))
	mockCognitoClient.AssertExpectations(t)
	mockCognitoClient.AssertNumberOfCalls(t, "CreateGroup", 0)
}
