package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"

	"github.com/panther-labs/panther/api/analysis/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

type getParams struct {
	ID        models.ID
	VersionID models.VersionID
}

// GetPolicy retrieves a policy from Dynamo or S3.
func GetPolicy(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	return handleGet(request, typePolicy)
}

// GetRule retrieves a rule from Dynamo or S3.
func GetRule(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	return handleGet(request, typeRule)
}

// Handle GET request for GetPolicy and GetRule
func handleGet(request *events.APIGatewayProxyRequest, codeType string) *events.APIGatewayProxyResponse {
	input, err := parseGet(request, codeType)
	if err != nil {
		return badRequest(err)
	}

	var item *tableItem
	if input.VersionID == "" {
		// Get latest version from Dynamo
		item, err = dynamoGet(input.ID, false)
	} else {
		// Get specific version from S3
		item, err = s3Get(input.ID, input.VersionID)
	}

	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}
	if item == nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}
	}

	if item.Type != codeType {
		// Item is the wrong type (e.g. a policy, not a rule)
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}
	}

	// Add current pass/fail information and convert to external Policy model
	if codeType == typePolicy {
		status, err := getComplianceStatus(input.ID)
		if err != nil {
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}
		return gatewayapi.MarshalResponse(item.Policy(status.Status), http.StatusOK)
	}

	return gatewayapi.MarshalResponse(item.Rule(), http.StatusOK)
}

// Parse GET parameters for GetPolicy and GetRule
func parseGet(request *events.APIGatewayProxyRequest, codeType string) (*getParams, error) {
	params := &getParams{
		VersionID: models.VersionID(request.QueryStringParameters["versionId"]),
	}

	idKey := "policyId"
	if codeType == typeRule {
		idKey = "ruleId"
	}
	id, err := url.QueryUnescape(request.QueryStringParameters[idKey])
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %s", idKey, err)
	}
	params.ID = models.ID(id)

	if err := params.ID.Validate(nil); err != nil {
		return nil, fmt.Errorf("invalid %s: %s", idKey, err)
	}

	if params.VersionID != "" {
		if err := params.VersionID.Validate(nil); err != nil {
			return nil, errors.New("invalid versionId: " + err.Error())
		}
	}

	return params, nil
}
