package outputs

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var (
	pagerDutyEndpoint  = "https://events.pagerduty.com/v2/enqueue"
	triggerEventAction = "trigger"
)

func pantherSeverityToPagerDuty(severity *string) (*string, *AlertDeliveryError) {
	switch *severity {
	case "INFO", "LOW":
		return aws.String("info"), nil
	case "MEDIUM":
		return aws.String("warning"), nil
	case "HIGH":
		return aws.String("error"), nil
	case "CRITICAL":
		return aws.String("critical"), nil
	default:
		return nil, &AlertDeliveryError{Message: "unknown severity" + aws.StringValue(severity)}
	}
}

// PagerDuty sends an alert to a pager duty integration endpoint.
func (client *OutputClient) PagerDuty(alert *alertmodels.Alert, config *outputmodels.PagerDutyConfig) *AlertDeliveryError {
	severity, err := pantherSeverityToPagerDuty(alert.Severity)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"summary":   *generateAlertTitle(alert),
		"severity":  aws.StringValue(severity),
		"timestamp": alert.CreatedAt.Format(time.RFC3339),
		"source":    "pantherlabs",
		"custom_details": map[string]string{
			"description": aws.StringValue(alert.PolicyDescription),
			"runbook":     aws.StringValue(alert.Runbook),
		},
	}

	pagerDutyRequest := map[string]interface{}{
		"payload":      payload,
		"routing_key":  *config.IntegrationKey,
		"event_action": triggerEventAction,
	}

	postInput := &PostInput{
		url:  &pagerDutyEndpoint,
		body: pagerDutyRequest,
	}

	return client.httpWrapper.post(postInput)
}
