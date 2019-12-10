package processor

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// Replace global logger with an in-memory observer for tests.
func mockLogger() *observer.ObservedLogs {
	core, mockLog := observer.New(zap.DebugLevel)
	zap.ReplaceGlobals(zap.New(core))
	return mockLog
}

type mockSns struct {
	mock.Mock
	snsiface.SNSAPI
}

func (m *mockSns) ConfirmSubscription(in *sns.ConfirmSubscriptionInput) (*sns.ConfirmSubscriptionOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*sns.ConfirmSubscriptionOutput), args.Error(1)
}

type mockSqs struct {
	mock.Mock
	sqsiface.SQSAPI
}

func (m *mockSqs) SendMessageBatch(in *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	args := m.Called(in)
	return args.Get(0).(*sqs.SendMessageBatchOutput), args.Error(1)
}
