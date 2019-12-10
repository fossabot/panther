package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/internal/compliance/aws_event_processor/processor"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

func handler(ctx context.Context, batch *events.SQSEvent) error {
	lambdalogger.ConfigureGlobal(ctx, nil)
	return processor.Handle(batch)
}

func main() {
	lambda.Start(handler)
}
