package delivery

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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/panther-labs/panther/internal/core/alert_delivery/outputs"
)

var (
	awsSession = session.Must(session.NewSession())

	// We will always need the Lambda client (to get output details)
	lambdaClient lambdaiface.LambdaAPI = lambda.New(awsSession)

	outputClient outputs.API = outputs.New(awsSession)

	// Lazy-load the SQS client - we only need it to retry failed alerts
	sqsClient sqsiface.SQSAPI
)

func getSQSClient() sqsiface.SQSAPI {
	if sqsClient == nil {
		sqsClient = sqs.New(awsSession)
	}
	return sqsClient
}
