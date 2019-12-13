package outputs

import (
	"github.com/aws/aws-sdk-go/aws"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var (
	opsgenieEndpoint = "https://api.opsgenie.com/v2/alerts"
)

var pantherToOpsGeniePriority = map[string]string{
	"CRITICAL": "P1",
	"HIGH":     "P2",
	"MEDIUM":   "P3",
	"LOW":      "P4",
	"INFO":     "P5",
}

// Opsgenie alert send an alert.
func (client *OutputClient) Opsgenie(
	alert *alertmodels.Alert, config *outputmodels.OpsgenieConfig) *AlertDeliveryError {

	tagsItem := aws.StringValueSlice(alert.Tags)

	description := "<strong>Description:</strong> " + aws.StringValue(alert.PolicyDescription)
	link := "\n<a href=\"" + generateURL(alert) + "\">Click here to view in the Panther UI</a>"
	runBook := "\n <strong>Runbook:</strong> " + aws.StringValue(alert.Runbook)
	severity := "\n <strong>Severity:</strong> " + aws.StringValue(alert.Severity)

	opsgenieRequest := map[string]interface{}{
		"message":     *generateAlertTitle(alert),
		"description": description + link + runBook + severity,
		"tags":        tagsItem,
		"priority":    pantherToOpsGeniePriority[aws.StringValue(alert.Severity)],
	}
	authorization := "GenieKey " + *config.APIKey
	accept := "application/json"
	requestHeader := map[string]*string{
		"Accept":        &accept,
		"Authorization": &authorization,
	}

	postInput := &PostInput{
		url:     &opsgenieEndpoint,
		body:    opsgenieRequest,
		headers: requestHeader,
	}
	return client.httpWrapper.post(postInput)
}
