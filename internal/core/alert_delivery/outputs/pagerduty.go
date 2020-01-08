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
	"time"

	"github.com/aws/aws-sdk-go/aws"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var (
	pagerDutyEndpoint  = "https://events.pagerduty.com/v2/enqueue"
	triggerEventAction = "trigger"
)

func pantherSeverityToPagerDuty(severity *string) (*string, *AlertDeliveryError) {
	switch *severity {
	case "INFO", "LOW":
		return aws.String("info"), nil
	case "MEDIUM":
		return aws.String("warning"), nil
	case "HIGH":
		return aws.String("error"), nil
	case "CRITICAL":
		return aws.String("critical"), nil
	default:
		return nil, &AlertDeliveryError{Message: "unknown severity" + aws.StringValue(severity)}
	}
}

// PagerDuty sends an alert to a pager duty integration endpoint.
func (client *OutputClient) PagerDuty(alert *alertmodels.Alert, config *outputmodels.PagerDutyConfig) *AlertDeliveryError {
	severity, err := pantherSeverityToPagerDuty(alert.Severity)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"summary":   *generateAlertTitle(alert),
		"severity":  aws.StringValue(severity),
		"timestamp": alert.CreatedAt.Format(time.RFC3339),
		"source":    "pantherlabs",
		"custom_details": map[string]string{
			"description": aws.StringValue(alert.PolicyDescription),
			"runbook":     aws.StringValue(alert.Runbook),
		},
	}

	pagerDutyRequest := map[string]interface{}{
		"payload":      payload,
		"routing_key":  *config.IntegrationKey,
		"event_action": triggerEventAction,
	}

	postInput := &PostInput{
		url:  &pagerDutyEndpoint,
		body: pagerDutyRequest,
	}

	return client.httpWrapper.post(postInput)
}
