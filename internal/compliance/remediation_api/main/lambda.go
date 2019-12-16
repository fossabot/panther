package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	apihandlers "github.com/panther-labs/panther/internal/compliance/remediation_api/handlers"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var methodHandlers = map[string]gatewayapi.RequestHandler{
	"GET /":                apihandlers.GetRemediations,
	"POST /remediate":      apihandlers.RemediateResource,
	"POST /remediateasync": apihandlers.RemediateResourceAsync,
}

func main() {
	lambda.Start(gatewayapi.LambdaProxy(methodHandlers))
}
