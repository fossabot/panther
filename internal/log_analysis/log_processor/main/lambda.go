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
	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/common"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/destinations"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/processor"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/sources"
	"github.com/panther-labs/panther/pkg/lambdalogger"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, event events.SQSEvent) error {
	lc, _ := lambdalogger.ConfigureGlobal(ctx, nil)
	return process(lc, event)
}

func process(lc *lambdacontext.LambdaContext, event events.SQSEvent) (err error) {
	operation := common.OpLogManager.Start(lc.InvokedFunctionArn, common.OpLogLambdaServiceDim)
	defer func() {
		operation.Stop()
		operation.Log(err, zap.Int("sqsMessageCount", len(event.Records)))
	}()

	// this is not likely to happen in production but needed to avoid opening sessions in tests w/no events
	if len(event.Records) == 0 {
		return err
	}

	dataStreams, err := sources.ReadSQSMessages(event.Records)
	if err != nil {
		return err
	}
	err = processor.Process(dataStreams, destinations.CreateDestination())
	return err
}
