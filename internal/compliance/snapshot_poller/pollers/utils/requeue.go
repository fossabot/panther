package utils

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
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
)

// SQS imposed maximum message delay
const MaxRequeueDelaySeconds = 900

var queueURL = os.Getenv("SNAPSHOT_QUEUE_URL")

// Requeue sends a scan request back to the poller input queue
func Requeue(scanRequest poller.ScanMsg, delay int64) {
	body, err := jsoniter.MarshalToString(scanRequest)
	if err != nil {
		zap.L().Error("unable to marshal requeue request", zap.Any("request", scanRequest))
		return
	}

	if delay > MaxRequeueDelaySeconds {
		delay = MaxRequeueDelaySeconds
	}

	sqsClient := sqs.New(session.Must(session.NewSession()))
	_, err = sqsClient.SendMessage(
		&sqs.SendMessageInput{
			MessageBody:  aws.String(body),
			QueueUrl:     &queueURL,
			DelaySeconds: aws.Int64(delay),
		})
	if err != nil {
		zap.L().Error("scan re-queueing failed", zap.Error(err), zap.Any("request", scanRequest))
	}
}
