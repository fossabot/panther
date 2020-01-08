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
	"strings"

	"github.com/aws/aws-sdk-go/aws"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

// MsTeams alert send an alert.
func (client *OutputClient) MsTeams(
	alert *alertmodels.Alert, config *outputmodels.MsTeamsConfig) *AlertDeliveryError {

	var tagsItem = aws.StringValueSlice(alert.Tags)

	link := "[Click here to view in the Panther UI](" + policyURLPrefix + aws.StringValue(alert.PolicyID) + ").\n"
	runBook := aws.StringValue(alert.Runbook)
	ruleDescription := aws.StringValue(alert.PolicyDescription)
	severity := aws.StringValue(alert.Severity)
	tags := strings.Join(tagsItem, ", ")

	msTeamsRequestBody := map[string]interface{}{
		"@context": "http://schema.org/extensions",
		"@type":    "MessageCard",
		"text":     *generateAlertTitle(alert),
		"sections": []interface{}{
			map[string]interface{}{
				"facts": []interface{}{
					map[string]string{"name": "Description", "value": ruleDescription},
					map[string]string{"name": "Runbook", "value": runBook},
					map[string]string{"name": "Severity", "value": severity},
					map[string]string{"name": "Tags", "value": tags},
				},
				"text": link,
			},
		},
		"potentialAction": []interface{}{
			map[string]interface{}{
				"@type": "OpenUri",
				"name":  "Click here to view in the Panther UI",
				"targets": []interface{}{
					map[string]string{
						"os":  "default",
						"uri": generateURL(alert),
					},
				},
			},
		},
	}

	accept := "application/json"
	requestHeader := map[string]*string{
		"Accept": &accept,
	}
	requestURL := *config.WebhookURL
	postInput := &PostInput{
		url:     &requestURL,
		body:    msTeamsRequestBody,
		headers: requestHeader,
	}
	return client.httpWrapper.post(postInput)
}
