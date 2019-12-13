package outputs

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var msTeamConfig = &outputmodels.MsTeamsConfig{
	WebhookURL: aws.String("msteam-url"),
}

func TestMsTeamsAlert(t *testing.T) {
	httpWrapper := &mockHTTPWrapper{}
	client := &OutputClient{httpWrapper: httpWrapper}

	var createdAtTime, _ = time.Parse(time.RFC3339, "2019-08-03T11:40:13Z")
	alert := &alertmodels.Alert{
		PolicyID:   aws.String("policyId"),
		CreatedAt:  &createdAtTime,
		OutputIDs:  aws.StringSlice([]string{"output-id"}),
		PolicyName: aws.String("policyName"),
		Severity:   aws.String("INFO"),
	}

	msTeamsPayload := map[string]interface{}{
		"@context": "http://schema.org/extensions",
		"@type":    "MessageCard",
		"text":     "Policy Failure: policyName",
		"sections": []interface{}{
			map[string]interface{}{
				"facts": []interface{}{
					map[string]string{"name": "Description", "value": ""},
					map[string]string{"name": "Runbook", "value": ""},
					map[string]string{"name": "Severity", "value": "INFO"},
					map[string]string{"name": "Tags", "value": ""},
				},
				"text": "[Click here to view in the Panther UI](https://panther.io/policies/policyId).\n",
			},
		},
		"potentialAction": []interface{}{
			map[string]interface{}{
				"@type": "OpenUri",
				"name":  "Click here to view in the Panther UI",
				"targets": []interface{}{
					map[string]string{
						"os":  "default",
						"uri": "https://panther.io/policies/policyId",
					},
				},
			},
		},
	}

	requestURL := *msTeamConfig.WebhookURL
	accept := "application/json"
	requestHeader := map[string]*string{
		"Accept": &accept,
	}
	expectedPostInput := &PostInput{
		url:     &requestURL,
		body:    msTeamsPayload,
		headers: requestHeader,
	}

	httpWrapper.On("post", expectedPostInput).Return((*AlertDeliveryError)(nil))

	require.Nil(t, client.MsTeams(alert, msTeamConfig))
	httpWrapper.AssertExpectations(t)
}
