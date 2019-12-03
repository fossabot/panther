package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// ModifyRule updates an existing rule.
func ModifyRule(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseUpdateRule(request)
	if err != nil {
		return badRequest(err)
	}

	item := &tableItem{
		Body:          input.Body,
		Description:   input.Description,
		DisplayName:   input.DisplayName,
		Enabled:       input.Enabled,
		ID:            input.ID,
		Reference:     input.Reference,
		ResourceTypes: input.LogTypes,
		Runbook:       input.Runbook,
		Severity:      input.Severity,
		Tags:          input.Tags,
		Tests:         input.Tests,
		Type:          typeRule,
	}

	if _, err := writeItem(item, input.UserID, aws.Bool(true)); err != nil {
		if err == errNotExists || err == errWrongType {
			// errWrongType means we tried to modify a rule which is actually a policy.
			// In this case return 404 - the rule you tried to modify does not exist.
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}
		}
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(item.Rule(), http.StatusOK)
}
