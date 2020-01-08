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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var jiraConfig = &outputmodels.JiraConfig{
	OrgDomain:  aws.String("https://panther-labs.atlassian.net"),
	ProjectKey: aws.String("QR"),
	UserName:   aws.String("username"),
	APIKey:     aws.String("apikey"),
	AssigneeID: aws.String("ae393k930390"),
}

func TestJiraAlert(t *testing.T) {
	httpWrapper := &mockHTTPWrapper{}
	client := &OutputClient{httpWrapper: httpWrapper}

	var createdAtTime, _ = time.Parse(time.RFC3339, "2019-08-03T11:40:13Z")
	alert := &alertmodels.Alert{
		PolicyID:          aws.String("ruleId"),
		CreatedAt:         &createdAtTime,
		OutputIDs:         aws.StringSlice([]string{"output-id"}),
		PolicyDescription: aws.String("policyDescription"),
		Severity:          aws.String("INFO"),
	}

	jiraPayload := map[string]interface{}{
		"fields": map[string]interface{}{
			"summary": "Policy Failure: ruleId",
			"description": "*Description:* policyDescription\n " +
				"[Click here to view in the Panther UI](https://panther.io/policies/ruleId)\n" +
				" *Runbook:* \n *Severity:* INFO\n *Tags:* ",
			"project": map[string]*string{
				"key": jiraConfig.ProjectKey,
			},
			"issuetype": map[string]string{
				"name": "Task",
			},
			"assignee": map[string]*string{
				"id": jiraConfig.AssigneeID,
			},
		},
	}
	auth := *jiraConfig.UserName + ":" + *jiraConfig.APIKey
	basicAuthToken := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	accept := "application/json"
	requestHeader := map[string]*string{
		"Authorization": &basicAuthToken,
		"Accept":        &accept,
	}
	requestEndpoint := "https://panther-labs.atlassian.net/rest/api/latest/issue/"
	expectedPostInput := &PostInput{
		url:     &requestEndpoint,
		body:    jiraPayload,
		headers: requestHeader,
	}

	httpWrapper.On("post", expectedPostInput).Return((*AlertDeliveryError)(nil))

	require.Nil(t, client.Jira(alert, jiraConfig))
	httpWrapper.AssertExpectations(t)
}
