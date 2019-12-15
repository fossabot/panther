package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/internal/core/organization_api/api"
	"github.com/panther-labs/panther/pkg/genericapi"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

var router = genericapi.NewRouter(nil, api.API{})

func lambdaHandler(ctx context.Context, input *models.LambdaInput) (interface{}, error) {
	lambdalogger.ConfigureGlobal(ctx, nil)
	return router.Handle(input)
}

func main() {
	lambda.Start(lambdaHandler)
}
