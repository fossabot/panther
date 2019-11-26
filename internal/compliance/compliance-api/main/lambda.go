package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"

	"github.com/panther-labs/panther/internal/compliance/compliance-api/handlers"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var methodHandlers = map[string]gatewayapi.RequestHandler{
	"GET /describe-org":      handlers.DescribeOrg,
	"GET /describe-policy":   handlers.DescribePolicy,
	"GET /describe-resource": handlers.DescribeResource,
	"GET /org-overview":      handlers.GetOrgOverview,
	"GET /status":            handlers.GetStatus,

	"POST /delete": handlers.DeleteStatus,
	"POST /status": handlers.SetStatus,
	"POST /update": handlers.UpdateMetadata,
}

func main() {
	envconfig.MustProcess("", &handlers.Env)
	lambda.Start(gatewayapi.LambdaProxy(methodHandlers))
}
