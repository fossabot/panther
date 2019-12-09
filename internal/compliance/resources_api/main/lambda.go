package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/internal/compliance/resources_api/handlers"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var methodHandlers = map[string]gatewayapi.RequestHandler{
	"POST /delete":      handlers.DeleteResources,
	"GET /list":         handlers.ListResources,
	"GET /org-overview": handlers.OrgOverview,
	"GET /resource":     handlers.GetResource,
	"POST /resource":    handlers.AddResources,
	"PATCH /resource":   handlers.ModifyResource,
}

func main() {
	handlers.Setup()
	lambda.Start(gatewayapi.LambdaProxy(methodHandlers))
}
