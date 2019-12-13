package outputs

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
