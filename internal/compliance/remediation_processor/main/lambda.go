package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/api/gateway/remediation/models"
	"github.com/panther-labs/panther/internal/compliance/remediation_api/remediation"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

var invoker = remediation.NewInvoker(session.Must(session.NewSession()))

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, event events.SQSEvent) error {
	lambdalogger.ConfigureGlobal(ctx, nil)

	for _, record := range event.Records {
		var input models.RemediateResource
		if err := jsoniter.UnmarshalFromString(record.Body, &input); err != nil {
			return err
		}
		if err := invoker.Remediate(&input); err != nil {
			return err
		}
	}
	return nil
}
