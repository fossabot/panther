package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	custommessage "github.com/panther-labs/panther/internal/core/custom_message/api"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

func lambdaHandler(ctx context.Context, event *events.CognitoEventUserPoolsCustomMessage) (
	*events.CognitoEventUserPoolsCustomMessage, error) {

	lambdalogger.ConfigureGlobal(ctx, nil)
	return custommessage.HandleEvent(event)
}

func main() {
	lambda.Start(lambdaHandler)
}
