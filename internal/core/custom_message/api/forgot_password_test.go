package custommessage

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
