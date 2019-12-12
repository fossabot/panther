package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/panther-labs/panther/internal/compliance/snapshot_scheduler/scheduler"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

func lambdaHandler(ctx context.Context, request events.CloudWatchEvent) error {
	lambdalogger.ConfigureGlobal(ctx, nil)
	return scheduler.PollAndIssueNewScans()
}

func main() {
	lambda.Start(lambdaHandler)
}
