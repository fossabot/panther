package gateway

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDescribeExecution(t *testing.T) {
	mockSvc := &MockSFN{}
	gw := &StepFunctionGateway{sfnClient: mockSvc}
	describeOut := &sfn.DescribeExecutionOutput{
		StartDate: &time.Time{},
		StopDate:  &time.Time{},
		Status:    aws.String("PASSED"),
	}
	mockSvc.
		On("DescribeExecution", mock.Anything).
		Return(describeOut, nil)

	result, err := gw.DescribeExecution(
		aws.String("fakeExecutionArn"),
	)
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestDescribeExecutionFailed(t *testing.T) {
	mockSvc := &MockSFN{}
	gw := &StepFunctionGateway{sfnClient: mockSvc}
	err := errors.New("sfn does not exist")
	mockSvc.
		On("DescribeExecution", mock.Anything).
		Return(&sfn.DescribeExecutionOutput{}, err)

	result, err := gw.DescribeExecution(
		aws.String("fakeExecutionArn"),
	)
	assert.Nil(t, result)
	assert.Error(t, err)
}
