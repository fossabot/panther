package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/api/gateway/resources/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// Convert a validation error into a 400 proxy response.
func badRequest(err error) *events.APIGatewayProxyResponse {
	errModel := &models.Error{Message: aws.String(err.Error())}
	return gatewayapi.MarshalResponse(errModel, http.StatusBadRequest)
}

func intMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}
