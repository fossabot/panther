package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// ModifyPolicy updates an existing policy.
func ModifyPolicy(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseUpdatePolicy(request)
	if err != nil {
		return badRequest(err)
	}

	item := &tableItem{
		AutoRemediationID:         input.AutoRemediationID,
		AutoRemediationParameters: input.AutoRemediationParameters,
		Body:                      input.Body,
		Description:               input.Description,
		DisplayName:               input.DisplayName,
		Enabled:                   input.Enabled,
		ID:                        input.ID,
		Reference:                 input.Reference,
		ResourceTypes:             input.ResourceTypes,
		Runbook:                   input.Runbook,
		Severity:                  input.Severity,
		Suppressions:              input.Suppressions,
		Tags:                      input.Tags,
		Tests:                     input.Tests,
		Type:                      typePolicy,
	}

	if _, err := writeItem(item, input.UserID, aws.Bool(true)); err != nil {
		if err == errNotExists || err == errWrongType {
			// errWrongType means we tried to modify a policy that is actually a rule.
			// In this case return 404 - the policy you tried to modify does not exist.
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}
		}
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	status, err := getComplianceStatus(input.ID)
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(item.Policy(status.Status), http.StatusOK)
}
