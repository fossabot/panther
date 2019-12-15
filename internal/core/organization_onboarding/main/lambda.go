package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/api/lambda/onboarding/models"
	"github.com/panther-labs/panther/internal/core/organization_onboarding/api"
	"github.com/panther-labs/panther/pkg/genericapi"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

var router = genericapi.NewRouter(nil, &api.API{})

func lambdaHandler(ctx context.Context, input *models.LambdaInput) (interface{}, error) {
	_, logger := lambdalogger.ConfigureGlobal(ctx, nil)
	defer func() { _ = logger.Sync() }()
	return router.Handle(input)
}

func main() {
	lambda.Start(lambdaHandler)
}
