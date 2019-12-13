package outputs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var githubConfig = &outputmodels.GithubConfig{RepoName: aws.String("profile/reponame"), Token: aws.String("github-token")}

func TestGithubAlert(t *testing.T) {
	httpWrapper := &mockHTTPWrapper{}
	client := &OutputClient{httpWrapper: httpWrapper}

	var createdAtTime, _ = time.Parse(time.RFC3339, "2019-08-03T11:40:13Z")
	alert := &alertmodels.Alert{
		PolicyID:          aws.String("ruleId"),
		CreatedAt:         &createdAtTime,
		OutputIDs:         aws.StringSlice([]string{"output-id"}),
		PolicyDescription: aws.String("description"),
		PolicyName:        aws.String("rule_name"),
		Severity:          aws.String("INFO"),
	}

	githubRequest := map[string]interface{}{
		"title": "Policy Failure: rule_name",
		"body": "**Description:** description\n " +
			"[Click here to view in the Panther UI](https://panther.io/policies/ruleId)\n" +
			" **Runbook:** \n **Severity:** INFO\n **Tags:** ",
	}

	authorization := "token " + *githubConfig.Token
	accept := "application/json"
	requestHeader := map[string]*string{
		"Authorization": &authorization,
		"Accept":        &accept,
	}
	requestEndpoint := "https://api.github.com/repos/profile/reponame/issues"
	expectedPostInput := &PostInput{
		url:     &requestEndpoint,
		body:    githubRequest,
		headers: requestHeader,
	}

	httpWrapper.On("post", expectedPostInput).Return((*AlertDeliveryError)(nil))

	require.Nil(t, client.Github(alert, githubConfig))
	httpWrapper.AssertExpectations(t)
}
