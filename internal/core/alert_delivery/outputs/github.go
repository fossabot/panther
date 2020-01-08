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

// Severity colors match those in the Panther UI
const (
	githubEndpoint = "https://api.github.com/repos/"
	requestType    = "/issues"
)

// Github alert send an issue.
func (client *OutputClient) Github(
	alert *alertmodels.Alert, config *outputmodels.GithubConfig) *AlertDeliveryError {

	var tagsItem = aws.StringValueSlice(alert.Tags)

	description := "**Description:** " + aws.StringValue(alert.PolicyDescription)
	link := "\n [Click here to view in the Panther UI](" + generateURL(alert) + ")"
	runBook := "\n **Runbook:** " + aws.StringValue(alert.Runbook)
	severity := "\n **Severity:** " + aws.StringValue(alert.Severity)
	tags := "\n **Tags:** " + strings.Join(tagsItem, ", ")

	githubRequest := map[string]interface{}{
		"title": aws.StringValue(generateAlertTitle(alert)),
		"body":  description + link + runBook + severity + tags,
	}

	accept := "application/json"
	token := "token " + *config.Token
	repoURL := githubEndpoint + *config.RepoName + requestType
	requestHeader := map[string]*string{
		"Authorization": &token,
		"Accept":        &accept,
	}

	postInput := &PostInput{
		url:     &repoURL,
		body:    githubRequest,
		headers: requestHeader,
	}
	return client.httpWrapper.post(postInput)
}
