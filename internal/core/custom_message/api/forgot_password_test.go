package custommessage

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/internal/core/users_api/gateway"
)

func TestHandleForgotPasswordGeneratePlainTextEmail(t *testing.T) {
	mockGateway := &gateway.MockUserGateway{}
	userGateway = mockGateway
	appDomainURL = "dev.runpanther.pizza"
	username := "user-123"
	poolID := "pool-123"
	codeParam := "123456"
	event := events.CognitoEventUserPoolsCustomMessage{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName:   username,
			UserPoolID: poolID,
		},
		Request: events.CognitoEventUserPoolsCustomMessageRequest{
			CodeParameter: codeParam,
		},
	}
	mockGateway.On("GetUser", &username, &poolID).Return(&models.User{
		GivenName:  aws.String("user-given-name-123"),
		FamilyName: aws.String("user-family-name-123"),
		Email:      aws.String("user@test.pizza"),
	}, nil)

	e, err := handleForgotPassword(&event)
	assert.Nil(t, err)
	assert.NotNil(t, e)
}
