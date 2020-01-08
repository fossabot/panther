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
	"encoding/base64"
	"strings"

	"github.com/aws/aws-sdk-go/aws"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

const (
	jiraEndpoint = "/rest/api/latest/issue/"
)

// Jira alert send an issue.
func (client *OutputClient) Jira(
	alert *alertmodels.Alert, config *outputmodels.JiraConfig) *AlertDeliveryError {

	var tagsItem = aws.StringValueSlice(alert.Tags)

	description := "*Description:* " + aws.StringValue(alert.PolicyDescription)
	link := "\n [Click here to view in the Panther UI](" + generateURL(alert) + ")"
	runBook := "\n *Runbook:* " + aws.StringValue(alert.Runbook)
	severity := "\n *Severity:* " + aws.StringValue(alert.Severity)
	tags := "\n *Tags:* " + strings.Join(tagsItem, ", ")

	fields := map[string]interface{}{
		"summary":     *generateAlertTitle(alert),
		"description": description + link + runBook + severity + tags,
		"project": map[string]*string{
			"key": config.ProjectKey,
		},
		"issuetype": map[string]string{
			"name": "Task",
		},
	}

	if config.AssigneeID != nil {
		fields["assignee"] = map[string]*string{
			"id": config.AssigneeID,
		}
	}

	jiraRequest := map[string]interface{}{
		"fields": fields,
	}

	auth := *config.UserName + ":" + *config.APIKey
	basicAuthToken := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	accept := "application/json"
	jiraRestURL := *config.OrgDomain + jiraEndpoint
	requestHeader := map[string]*string{
		"Accept":        &accept,
		"Authorization": &basicAuthToken,
	}

	postInput := &PostInput{
		url:     &jiraRestURL,
		body:    jiraRequest,
		headers: requestHeader,
	}
	return client.httpWrapper.post(postInput)
}
