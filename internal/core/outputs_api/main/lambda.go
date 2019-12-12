package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/internal/core/outputs_api/api"
	"github.com/panther-labs/panther/pkg/genericapi"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

var router *genericapi.Router

func init() {
	validator, err := models.Validator()
	if err != nil {
		panic(err)
	}
	router = genericapi.NewRouter(validator, api.API{})
}

func lambdaHandler(ctx context.Context, input *models.LambdaInput) (interface{}, error) {
	lambdalogger.ConfigureGlobal(ctx, nil)
	return router.Handle(input)
}

func main() {
	lambda.Start(lambdaHandler)
}
