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
	"github.com/aws/aws-sdk-go/aws"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var (
	opsgenieEndpoint = "https://api.opsgenie.com/v2/alerts"
)

var pantherToOpsGeniePriority = map[string]string{
	"CRITICAL": "P1",
	"HIGH":     "P2",
	"MEDIUM":   "P3",
	"LOW":      "P4",
	"INFO":     "P5",
}

// Opsgenie alert send an alert.
func (client *OutputClient) Opsgenie(
	alert *alertmodels.Alert, config *outputmodels.OpsgenieConfig) *AlertDeliveryError {

	tagsItem := aws.StringValueSlice(alert.Tags)

	description := "<strong>Description:</strong> " + aws.StringValue(alert.PolicyDescription)
	link := "\n<a href=\"" + generateURL(alert) + "\">Click here to view in the Panther UI</a>"
	runBook := "\n <strong>Runbook:</strong> " + aws.StringValue(alert.Runbook)
	severity := "\n <strong>Severity:</strong> " + aws.StringValue(alert.Severity)

	opsgenieRequest := map[string]interface{}{
		"message":     *generateAlertTitle(alert),
		"description": description + link + runBook + severity,
		"tags":        tagsItem,
		"priority":    pantherToOpsGeniePriority[aws.StringValue(alert.Severity)],
	}
	authorization := "GenieKey " + *config.APIKey
	accept := "application/json"
	requestHeader := map[string]*string{
		"Accept":        &accept,
		"Authorization": &authorization,
	}

	postInput := &PostInput{
		url:     &opsgenieEndpoint,
		body:    opsgenieRequest,
		headers: requestHeader,
	}
	return client.httpWrapper.post(postInput)
}
