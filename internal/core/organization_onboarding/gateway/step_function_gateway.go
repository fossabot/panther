package gateway

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	sfnI "github.com/aws/aws-sdk-go/service/sfn/sfniface"
	"github.com/stretchr/testify/mock"

	"github.com/panther-labs/panther/api/lambda/onboarding/models"
)

// API defines the interface for the user gateway which can be used for mocking.
type API interface {
	DescribeExecution(executionArn *string) (*models.GetOnboardingStatusOutput, error)
}

// StepFunctionGateway encapsulates a service to AWS Step Function.
type StepFunctionGateway struct {
	sfnClient sfnI.SFNAPI
}

// The StepFunctionGateway must satisfy the API interface.
var _ API = (*StepFunctionGateway)(nil)

// New creates a new StepFunctionClient client.
func New(sess *session.Session) *StepFunctionGateway {
	return &StepFunctionGateway{
		sfnClient: sfn.New(sess),
	}
}

// MockSFN is a mock CloudTrail client.
type MockSFN struct {
	sfnI.SFNAPI
	mock.Mock
}

// DescribeExecution is a mock function to return fake CloudTrail data.
func (m *MockSFN) DescribeExecution(
	in *sfn.DescribeExecutionInput,
) (*sfn.DescribeExecutionOutput, error) {

	args := m.Called(in)
	return args.Get(0).(*sfn.DescribeExecutionOutput), args.Error(1)
}
