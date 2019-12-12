package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// CreateRule adds a new rule to the Dynamo table.
func CreateRule(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
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

	if _, err := writeItem(item, input.UserID, aws.Bool(false)); err != nil {
		if err == errExists {
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusConflict}
		}
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(item.Rule(), http.StatusCreated)
}

// body parsing shared by CreateRule and ModifyRule
func parseUpdateRule(request *events.APIGatewayProxyRequest) (*models.UpdateRule, error) {
	var result models.UpdateRule
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	if err := result.Validate(nil); err != nil {
		return nil, err
	}

	return &result, nil
}
