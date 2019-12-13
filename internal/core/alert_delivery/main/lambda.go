package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/panther-labs/panther/internal/core/alert_delivery/delivery"
	"github.com/panther-labs/panther/internal/core/alert_delivery/models"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

var validate = validator.New()

func lambdaHandler(ctx context.Context, event events.SQSEvent) {
	_, logger := lambdalogger.ConfigureGlobal(ctx, nil)
	var alerts []*models.Alert

	for _, record := range event.Records {
		alert := &models.Alert{}
		if err := jsoniter.UnmarshalFromString(record.Body, alert); err != nil {
			logger.Error("failed to parse SQS message", zap.Error(err))
			continue
		}
		if err := validate.Struct(alert); err != nil {
			logger.Error("invalid message received", zap.Error(err))
			continue
		}
		alerts = append(alerts, alert)
	}

	if len(alerts) > 0 {
		delivery.HandleAlerts(alerts)
	} else {
		logger.Info("no alerts to process")
	}
}

func main() {
	lambda.Start(lambdaHandler)
}
