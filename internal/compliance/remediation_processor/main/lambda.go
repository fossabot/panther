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
