package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/api"
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

func lambdaHandler(ctx context.Context, request *models.LambdaInput) (interface{}, error) {
	lambdalogger.ConfigureGlobal(ctx, nil)
	return router.Handle(request)
}

func main() {
	lambda.Start(lambdaHandler)
}
