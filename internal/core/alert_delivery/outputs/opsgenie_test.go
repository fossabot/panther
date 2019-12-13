package outputs

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	outputmodels "github.com/panther-labs/panther/api/lambda/outputs/models"
	alertmodels "github.com/panther-labs/panther/internal/core/alert_delivery/models"
)

var opsgenieConfig = &outputmodels.OpsgenieConfig{APIKey: aws.String("apikey")}

func TestOpsgenieAlert(t *testing.T) {
	httpWrapper := &mockHTTPWrapper{}
	client := &OutputClient{httpWrapper: httpWrapper}

	var createdAtTime, _ = time.Parse(time.RFC3339, "2019-08-03T11:40:13Z")
	alert := &alertmodels.Alert{
		PolicyID:   aws.String("policyId"),
		CreatedAt:  &createdAtTime,
		OutputIDs:  aws.StringSlice([]string{"output-id"}),
		PolicyName: aws.String("policyName"),
		Severity:   aws.String("CRITICAL"),
	}

	opsgenieRequest := map[string]interface{}{
		"message": "Policy Failure: policyName",
		"description": strings.Join([]string{
			"<strong>Description:</strong> ",
			"<a href=\"https://panther.io/policies/policyId\">Click here to view in the Panther UI</a>",
			" <strong>Runbook:</strong> ",
			" <strong>Severity:</strong> CRITICAL",
		}, "\n"),
		"tags":     []string{},
		"priority": "P1",
	}

	authorization := "GenieKey " + *opsgenieConfig.APIKey

	accept := "application/json"
	requestHeader := map[string]*string{
		"Accept":        &accept,
		"Authorization": &authorization,
	}
	requestEndpoint := "https://api.opsgenie.com/v2/alerts"
	expectedPostInput := &PostInput{
		url:     &requestEndpoint,
		body:    opsgenieRequest,
		headers: requestHeader,
	}

	httpWrapper.On("post", expectedPostInput).Return((*AlertDeliveryError)(nil))

	require.Nil(t, client.Opsgenie(alert, opsgenieConfig))
	httpWrapper.AssertExpectations(t)
}
