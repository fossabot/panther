package gateway

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"
)

var testAddUserToGroupInput = &provider.AdminAddUserToGroupInput{
	GroupName:  aws.String("Admin"),
	Username:   aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766"),
	UserPoolId: aws.String("us-west-2_ZlG7Ldp1K"),
}

func TestAddUserToGroup(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On(
		"AdminAddUserToGroup", testAddUserToGroupInput).Return(&provider.AdminAddUserToGroupOutput{}, nil)

	assert.NoError(t, gw.AddUserToGroup(
		testAddUserToGroupInput.Username,
		testAddUserToGroupInput.GroupName,
		testAddUserToGroupInput.UserPoolId,
	))
	mockCognitoClient.AssertExpectations(t)
}

func TestAddUserToGroupFailure(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On("AdminAddUserToGroup", testAddUserToGroupInput).Return(
		&provider.AdminAddUserToGroupOutput{}, &genericapi.AWSError{})

	assert.Error(t, gw.AddUserToGroup(
		testAddUserToGroupInput.Username,
		testAddUserToGroupInput.GroupName,
		testAddUserToGroupInput.UserPoolId,
	))
	mockCognitoClient.AssertExpectations(t)
}
