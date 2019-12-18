package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/panther-labs/panther/pkg/lambdalogger"

	"github.com/panther-labs/panther/internal/log_analysis/alert_merger/merger"
)

var validate = validator.New()

func main() {
	lambda.Start(Handler)
}

// Handler is the entry point for the alert merger Lambda
func Handler(ctx context.Context, event events.SQSEvent) error {
	_, logger := lambdalogger.ConfigureGlobal(ctx, nil)
	for _, record := range event.Records {
		input := &merger.AlertNotification{}
		if err := jsoniter.UnmarshalFromString(record.Body, input); err != nil {
			logger.Warn("failed to unmarshall event", zap.Error(err))
			return err
		}

		if err := validate.Struct(input); err != nil {
			logger.Error("invalid message received", zap.Error(err))
			continue
		}

		if err := merger.Handle(input); err != nil {
			logger.Warn("encountered issue while processing event", zap.Error(err))
			return err
		}
	}
	return nil
}
