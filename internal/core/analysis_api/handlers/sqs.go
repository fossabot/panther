package handlers

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
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

// Queue a policy for re-analysis (evaluate against all applicable resources).
//
// This ensures policy changes are reflected almost immediately (instead of waiting for daily scan).
func queuePolicy(policy *tableItem) error {
	body, err := jsoniter.MarshalToString(policy.Policy(""))
	if err != nil {
		zap.L().Error("failed to marshal policy", zap.Error(err))
		return err
	}

	zap.L().Info("queueing policy for analysis",
		zap.String("policyId", string(policy.ID)),
		zap.String("resourceQueueURL", env.ResourceQueueURL))
	_, err = sqsClient.SendMessage(
		&sqs.SendMessageInput{MessageBody: &body, QueueUrl: &env.ResourceQueueURL})
	return err
}
