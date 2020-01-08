package main

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
