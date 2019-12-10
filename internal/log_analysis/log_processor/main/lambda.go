package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/processor"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/sources"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, event events.SQSEvent) error {
	lambdalogger.ConfigureGlobal(ctx, nil)

	zap.L().Info("num_sqs_messages", zap.Int("count", len(event.Records)))
	messages := make([]*string, len(event.Records))
	for i, message := range event.Records {
		messages[i] = aws.String(message.Body)
	}

	buffers, err := sources.ReadData(messages)
	if err != nil {
		return err
	}
	return processor.Handle(buffers)
}
