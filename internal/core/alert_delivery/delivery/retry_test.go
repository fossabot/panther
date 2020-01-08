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
	"errors"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type mockSQSClient struct {
	sqsiface.SQSAPI
	err bool
}

var sqsMessages int // store number of messages here for tests to verify

func (m mockSQSClient) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	if m.err {
		return nil, errors.New("internal service error")
	}
	sqsMessages = len(input.Entries)
	return &sqs.SendMessageBatchOutput{
		Successful: make([]*sqs.SendMessageBatchResultEntry, len(input.Entries)),
	}, nil
}
