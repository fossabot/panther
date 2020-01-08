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
