package outputs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var slackConfig = &outputmodels.SlackConfig{WebhookURL: aws.String("slack-channel-url")}

func TestSlackAlert(t *testing.T) {
	httpWrapper := &mockHTTPWrapper{}
	client := &OutputClient{httpWrapper: httpWrapper}

	createdAtTime := time.Now()
	alert := &alertmodels.Alert{
		PolicyID:   aws.String("policyId"),
		CreatedAt:  &createdAtTime,
		OutputIDs:  aws.StringSlice([]string{"output-id"}),
		PolicyName: aws.String("policyName"),
		Severity:   aws.String("INFO"),
	}

	expectedPostPayload := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{"color": "#47b881",
				"fallback": "Policy Failure: policyName",
				"fields": []map[string]interface{}{
					{
						"short": false,
						"value": "<https://panther.io/policies/policyId|Click here to view in the Panther UI>",
					},
					{
						"short": false,
						"title": "Runbook",
						"value": "",
					},
					{
						"short": true,
						"title": "Severity",
						"value": "INFO",
					},
				},
				"title": "Policy Failure: policyName",
			},
		},
	}
	requestEndpoint := "slack-channel-url"
	expectedPostInput := &PostInput{
		url:  &requestEndpoint,
		body: expectedPostPayload,
	}

	httpWrapper.On("post", expectedPostInput).Return((*AlertDeliveryError)(nil))

	require.Nil(t, client.Slack(alert, slackConfig))
	httpWrapper.AssertExpectations(t)
}
