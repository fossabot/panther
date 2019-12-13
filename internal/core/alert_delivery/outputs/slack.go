package outputs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

// Severity colors match those in the Panther UI
var severityColors = map[string]string{
	"CRITICAL": "#425a70",
	"HIGH":     "#cb2e2e",
	"MEDIUM":   "#d9822b",
	"LOW":      "#f7d154",
	"INFO":     "#47b881",
}

// Slack sends an alert to a slack channel.
func (client *OutputClient) Slack(alert *alertmodels.Alert, config *outputmodels.SlackConfig) *AlertDeliveryError {
	messageField := fmt.Sprintf("<%s|%s>",
		generateURL(alert),
		"Click here to view in the Panther UI")
	fields := []map[string]interface{}{
		{
			"value": messageField,
			"short": false,
		},
		{
			"title": "Runbook",
			"value": aws.StringValue(alert.Runbook),
			"short": false,
		},
		{
			"title": "Severity",
			"value": aws.StringValue(alert.Severity),
			"short": true,
		},
	}

	payload := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"fallback": aws.StringValue(generateAlertTitle(alert)),
				"color":    severityColors[aws.StringValue(alert.Severity)],
				"title":    aws.StringValue(generateAlertTitle(alert)),
				"fields":   fields,
			},
		},
	}
	requestEndpoint := *config.WebhookURL
	postInput := &PostInput{
		url:  &requestEndpoint,
		body: payload,
	}

	return client.httpWrapper.post(postInput)
}
