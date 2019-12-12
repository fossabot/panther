package models

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

var mockID = aws.String("39c78f99-0149-482e-8f7b-090d75f253bf")

func TestUpdateUserNoFields(t *testing.T) {
	assert.NotNil(t, Validator().Struct(&UpdateUserInput{ID: mockID}))
}

func TestUpdateUserBlankField(t *testing.T) {
	assert.Error(t, Validator().Struct(&UpdateUserInput{
		ID:         mockID,
		GivenName:  aws.String("given-name"),
		Role:       aws.String(""),
		UserPoolID: aws.String("fakePoolId"),
	}))
}

func TestUpdateUserOneField(t *testing.T) {
	assert.NoError(t, Validator().Struct(&UpdateUserInput{
		ID:         mockID,
		Role:       aws.String("Admin"),
		UserPoolID: aws.String("fakePoolId"),
	}))
}

func TestUpdateUserAllFields(t *testing.T) {
	assert.NoError(t, Validator().Struct(&UpdateUserInput{
		ID:          mockID,
		GivenName:   aws.String("given-name"),
		FamilyName:  aws.String("family-name"),
		PhoneNumber: aws.String("phone-num"),
		Role:        aws.String("Admin"),
		UserPoolID:  aws.String("fakePoolId"),
	}))
}
