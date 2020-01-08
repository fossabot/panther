package gateway

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
