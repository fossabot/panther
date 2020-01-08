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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var createdAtTime, _ = time.Parse(time.RFC3339, "2019-05-03T11:40:13Z")

var pagerDutyAlert = &alertmodels.Alert{
	PolicyName: aws.String("policyName"),
	PolicyID:   aws.String("policyId"),
	Severity:   aws.String("INFO"),
	Runbook:    aws.String("runbook"),
	CreatedAt:  &createdAtTime,
}
var pagerDutyConfig = &outputmodels.PagerDutyConfig{
	IntegrationKey: aws.String("integrationKey"),
}

func TestSendPagerDutyAlert(t *testing.T) {
	httpWrapper := &mockHTTPWrapper{}
	outputClient := &OutputClient{httpWrapper: httpWrapper}

	expectedPostPayload := map[string]interface{}{
		"event_action": "trigger",
		"payload": map[string]interface{}{
			"custom_details": map[string]string{
				"description": "",
				"runbook":     "runbook",
			},
			"severity":  "info",
			"source":    "pantherlabs",
			"summary":   "Policy Failure: policyName",
			"timestamp": "2019-05-03T11:40:13Z",
		},
		"routing_key": "integrationKey",
	}
	requestEndpoint := "https://events.pagerduty.com/v2/enqueue"
	expectedPostInput := &PostInput{
		url:  &requestEndpoint,
		body: expectedPostPayload,
	}

	httpWrapper.On("post", expectedPostInput).Return((*AlertDeliveryError)(nil))
	result := outputClient.PagerDuty(pagerDutyAlert, pagerDutyConfig)

	assert.Nil(t, result)
	httpWrapper.AssertExpectations(t)
}

func TestSendPagerDutyAlertPostError(t *testing.T) {
	httpWrapper := &mockHTTPWrapper{}
	outputClient := &OutputClient{httpWrapper: httpWrapper}

	httpWrapper.On("post", mock.Anything).Return(&AlertDeliveryError{Message: "Exception"})

	require.Error(t, outputClient.PagerDuty(pagerDutyAlert, pagerDutyConfig))
	httpWrapper.AssertExpectations(t)
}
