package delivery

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
