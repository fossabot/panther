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
