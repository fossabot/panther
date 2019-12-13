package outputs

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
