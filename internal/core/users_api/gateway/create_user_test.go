package gateway

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	provider "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/pkg/genericapi"
)

var testCreateUserInput = &CreateUserInput{
	GivenName:   aws.String("Joe"),
	FamilyName:  aws.String("Blow"),
	Email:       aws.String("joe.blow@toe.com"),
	PhoneNumber: aws.String("+11234567890"),
	UserPoolID:  aws.String("userPoolId"),
}

var testAdminCreateUserInput = &provider.AdminCreateUserInput{
	DesiredDeliveryMediums: []*string{aws.String("EMAIL")},
	ForceAliasCreation:     aws.Bool(false),
	UserAttributes: []*provider.AttributeType{
		{
			Name:  aws.String("given_name"),
			Value: testCreateUserInput.GivenName,
		},
		{
			Name:  aws.String("family_name"),
			Value: testCreateUserInput.FamilyName,
		},
		{
			Name:  aws.String("email"),
			Value: testCreateUserInput.Email,
		},
		{
			Name:  aws.String("phone_number"),
			Value: testCreateUserInput.PhoneNumber,
		},
		{
			Name:  aws.String("email_verified"),
			Value: aws.String("true"),
		},
	},
	Username:   testCreateUserInput.Email,
	UserPoolId: testCreateUserInput.UserPoolID,
}

func TestCreateUser(t *testing.T) {
	testUserID := aws.String("bc010600-b2d6-4a8d-92ac-d4f8bd209766")
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On(
		"AdminCreateUser", testAdminCreateUserInput).Return(&provider.AdminCreateUserOutput{
		User: &provider.UserType{
			Username: testUserID,
		},
	}, nil)

	id, err := gw.CreateUser(testCreateUserInput)

	assert.Equal(t, id, testUserID)
	assert.NoError(t, err)
	mockCognitoClient.AssertExpectations(t)
}

func TestCreateUserFailed(t *testing.T) {
	mockIamService := &MockIamService{}
	mockCognitoClient := &MockCognitoClient{}
	gw := &UsersGateway{userPoolClient: mockCognitoClient, iamService: mockIamService}

	mockCognitoClient.On("AdminCreateUser", testAdminCreateUserInput).Return(
		&provider.AdminCreateUserOutput{}, &genericapi.AWSError{})

	id, err := gw.CreateUser(testCreateUserInput)

	assert.Nil(t, id)
	assert.Error(t, err)
	mockCognitoClient.AssertExpectations(t)
}
