package outputs

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
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

type mockSesClient struct {
	sesiface.SESAPI
	mock.Mock
}

var alert = &alertmodels.Alert{
	PolicyName:        aws.String("policyName"),
	PolicyID:          aws.String("policyId"),
	PolicyDescription: aws.String("policyDescription"),
	Severity:          aws.String("severity"),
	Runbook:           aws.String("runbook"),
}
var outputConfig = &outputmodels.EmailConfig{
	DestinationAddress: aws.String("destinationAddress"),
}

func (m *mockSesClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ses.SendEmailOutput), args.Error(1)
}

func init() {
	policyURLPrefix = "https://panther.io/policies/"
	alertURLPrefix = "https://panther.io/alerts/"
	sesConfigurationSet = "sesConfigurationSet"
}

func TestSendEmail(t *testing.T) {
	client := &mockSesClient{}
	outputClient := &OutputClient{sesClient: client, mailFrom: aws.String("email@email.com")}

	expectedEmailInput := &ses.SendEmailInput{
		ConfigurationSetName: aws.String("sesConfigurationSet"),
		Source:               aws.String("email@email.com"),
		Destination:          &ses.Destination{ToAddresses: []*string{aws.String("destinationAddress")}},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Policy Failure: policyName"),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data: aws.String("<h2>Message</h2><a href='https://panther.io/policies/policyId'>" +
						"policyName failed on new resources</a><br><h2>Severity</h2>severity<br>" +
						"<h2>Runbook</h2>runbook<br><h2>Description</h2>policyDescription"),
				},
			},
		},
	}

	client.On("SendEmail", expectedEmailInput).Return(&ses.SendEmailOutput{}, nil)
	result := outputClient.Email(alert, outputConfig)
	assert.Nil(t, result)
	client.AssertExpectations(t)
}

func TestSendEmailRule(t *testing.T) {
	client := &mockSesClient{}
	outputClient := &OutputClient{sesClient: client, mailFrom: aws.String("email@email.com")}

	var alert = &alertmodels.Alert{
		PolicyName:        aws.String("ruleName"),
		PolicyID:          aws.String("ruleId"),
		PolicyDescription: aws.String("ruleDescription"),
		Severity:          aws.String("severity"),
		Runbook:           aws.String("runbook"),
		Type:              aws.String(alertmodels.RuleType),
		AlertID:           aws.String("alertId"),
	}

	expectedEmailInput := &ses.SendEmailInput{
		ConfigurationSetName: aws.String("sesConfigurationSet"),
		Source:               aws.String("email@email.com"),
		Destination:          &ses.Destination{ToAddresses: []*string{aws.String("destinationAddress")}},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("New Alert: ruleName"),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data: aws.String("<h2>Message</h2><a href='https://panther.io/alerts/alertId'>" +
						"ruleName failed</a><br><h2>Severity</h2>severity<br>" +
						"<h2>Runbook</h2>runbook<br><h2>Description</h2>ruleDescription"),
				},
			},
		},
	}

	client.On("SendEmail", expectedEmailInput).Return(&ses.SendEmailOutput{}, nil)
	result := outputClient.Email(alert, outputConfig)
	assert.Nil(t, result)
	client.AssertExpectations(t)
}

func TestSendEmailPermanentError(t *testing.T) {
	client := &mockSesClient{}
	outputClient := &OutputClient{sesClient: client}

	client.On("SendEmail", mock.Anything).Return(
		&ses.SendEmailOutput{},
		errors.New("message failed"),
	)

	result := outputClient.Email(alert, outputConfig)
	require.Error(t, result)
	assert.Equal(t, &AlertDeliveryError{Message: "request failed message failed", Permanent: true}, result)
	client.AssertExpectations(t)
}

func TestSendEmailTemporaryError(t *testing.T) {
	client := &mockSesClient{}
	outputClient := &OutputClient{sesClient: client}

	client.On("SendEmail", mock.Anything).Return(
		&ses.SendEmailOutput{},
		awserr.New(ses.ErrCodeMessageRejected, "Message rejected", nil),
	)

	result := outputClient.Email(alert, outputConfig)
	require.Error(t, result)
	assert.Equal(t, &AlertDeliveryError{Message: "request failed MessageRejected: Message rejected", Permanent: false}, result)
	client.AssertExpectations(t)
}
