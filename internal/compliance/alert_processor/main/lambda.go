package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/panther-labs/panther/internal/compliance/alert_processor/models"
	"github.com/panther-labs/panther/internal/compliance/alert_processor/processor"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

var validate = validator.New()

func main() {
	lambda.Start(reporterHandler)
}

func reporterHandler(ctx context.Context, event events.SQSEvent) error {
	_, logger := lambdalogger.ConfigureGlobal(ctx, nil)
	for _, record := range event.Records {
		var input models.ComplianceNotification
		if err := jsoniter.UnmarshalFromString(record.Body, &input); err != nil {
			zap.L().Warn("failed to unmarshall event", zap.Error(err))
			return err
		}

		if err := validate.Struct(input); err != nil {
			logger.Error("invalid message received", zap.Error(err))
			continue
		}

		if err := processor.Handle(&input); err != nil {
			zap.L().Warn("encountered issue while processing event", zap.Error(err))
			return err
		}
	}
	return nil
}
