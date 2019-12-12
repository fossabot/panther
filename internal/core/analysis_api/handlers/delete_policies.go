package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
)

// DeletePolicies marks one or more policies as deleted.
func DeletePolicies(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseDeletePolicies(request)
	if err != nil {
		return badRequest(err)
	}

	if err = dynamoBatchDelete(input); err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	if err = s3BatchDelete(input); err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	if err = complianceBatchDelete(input.Policies, []string{}); err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

func parseDeletePolicies(request *events.APIGatewayProxyRequest) (*models.DeletePolicies, error) {
	var result models.DeletePolicies
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	if err := result.Validate(nil); err != nil {
		return nil, err
	}

	return &result, nil
}
