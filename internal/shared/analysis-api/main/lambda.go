package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/internal/shared/analysis-api/handlers"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var methodHandlers = map[string]gatewayapi.RequestHandler{
	// Policies only
	"GET /list":      handlers.ListPolicies,
	"GET /policy":    handlers.GetPolicy,
	"POST /policy":   handlers.CreatePolicy,
	"POST /suppress": handlers.Suppress,
	"POST /update":   handlers.ModifyPolicy,
	"POST /upload":   handlers.BulkUpload,

	// Rules only
	"GET /rule":         handlers.GetRule,
	"POST /rule":        handlers.CreateRule,
	"GET /rule/list":    handlers.ListRules,
	"POST /rule/update": handlers.ModifyRule,

	// Rules and Policies
	"POST /delete": handlers.DeletePolicies,
	"GET /enabled": handlers.GetEnabledPolicies,
	"POST /test":   handlers.TestPolicy,
}

func main() {
	handlers.Setup()
	lambda.Start(gatewayapi.LambdaProxy(methodHandlers))
}
