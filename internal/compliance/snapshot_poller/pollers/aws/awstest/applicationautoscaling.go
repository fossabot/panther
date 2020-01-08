package awstest

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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/stretchr/testify/mock"
)

// Example ApplicationAutoScaling return values
var (
	ExampleDescribeScalableTargetsOutput = &applicationautoscaling.DescribeScalableTargetsOutput{
		ScalableTargets: []*applicationautoscaling.ScalableTarget{
			{
				ServiceNamespace:  aws.String("dynamodb"),
				ResourceId:        aws.String("table/example-table"),
				ScalableDimension: aws.String("dynamodb:table:ReadCapacityUnits"),
				MinCapacity:       aws.Int64(5),
				MaxCapacity:       aws.Int64(4000),
				RoleARN:           aws.String("arn:aws:iam::123456789012:role/aws-service-role/dynamodb.application-autoscaling.amazonaws.com/AutoScalingRole"),
				CreationTime:      ExampleDate,
			},
			{
				ServiceNamespace:  aws.String("dynamodb"),
				ResourceId:        aws.String("table/example-table"),
				ScalableDimension: aws.String("dynamodb:table:WriteCapacityUnits"),
				MinCapacity:       aws.Int64(5),
				MaxCapacity:       aws.Int64(40000),
				RoleARN:           aws.String("arn:aws:iam::123456789012:role/aws-service-role/dynamodb.application-autoscaling.amazonaws.com/AutoScalingRole"),
				CreationTime:      ExampleDate,
			},
		},
	}

	svcApplicationAutoScalingSetupCalls = map[string]func(*MockApplicationAutoScaling){
		"DescribeScalableTargetsPages": func(svc *MockApplicationAutoScaling) {
			svc.On("DescribeScalableTargetsPages", mock.Anything).
				Return(nil)
		},
	}

	svcApplicationAutoScalingSetupCallsError = map[string]func(*MockApplicationAutoScaling){
		"DescribeScalableTargetsPages": func(svc *MockApplicationAutoScaling) {
			svc.On("DescribeScalableTargetsPages", mock.Anything).
				Return(errors.New("ApplicationAutoScaling.DescribeScalableTargetsPages"))
		},
	}

	MockApplicationAutoScalingForSetup = &MockApplicationAutoScaling{}
)

// Application Auto Scaling mock

// SetupMockApplicationAutoScaling is used to override the Application Auto Scaling Client initializer
func SetupMockApplicationAutoScaling(sess *session.Session, cfg *aws.Config) interface{} {
	return MockApplicationAutoScalingForSetup
}

// MockApplicationAutoScaling is a mock Application Auto Scaling client
type MockApplicationAutoScaling struct {
	applicationautoscalingiface.ApplicationAutoScalingAPI
	mock.Mock
}

// BuildMockApplicationAutoScalingSvc builds and returns a MockApplicationAutoScaling struct
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockApplicationAutoScalingSvc(funcs []string) (mockSvc *MockApplicationAutoScaling) {
	mockSvc = &MockApplicationAutoScaling{}
	for _, f := range funcs {
		svcApplicationAutoScalingSetupCalls[f](mockSvc)
	}
	return
}

// BuildMockApplicationAutoScalingSvcError builds and returns a MockApplicationAutoScaling struct with errors set
//
// Additionally, the appropriate calls to On and Return are made based on the strings passed in
func BuildMockApplicationAutoScalingSvcError(funcs []string) (mockSvc *MockApplicationAutoScaling) {
	mockSvc = &MockApplicationAutoScaling{}
	for _, f := range funcs {
		svcApplicationAutoScalingSetupCallsError[f](mockSvc)
	}
	return
}

// BuildApplicationAutoScalingServiceSvcAll builds and returns a MockApplicationAutoScaling struct
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockApplicationAutoScalingSvcAll() (mockSvc *MockApplicationAutoScaling) {
	mockSvc = &MockApplicationAutoScaling{}
	for _, f := range svcApplicationAutoScalingSetupCalls {
		f(mockSvc)
	}
	return
}

// BuildMockApplicationAutoScalingSvcAllError builds and returns a MockApplicationAutoScaling struct with errors set
//
// Additionally, the appropriate calls to On and Return are made for all possible function calls
func BuildMockApplicationAutoScalingSvcAllError() (mockSvc *MockApplicationAutoScaling) {
	mockSvc = &MockApplicationAutoScaling{}
	for _, f := range svcApplicationAutoScalingSetupCallsError {
		f(mockSvc)
	}
	return
}

func (m *MockApplicationAutoScaling) DescribeScalableTargetsPages(
	in *applicationautoscaling.DescribeScalableTargetsInput,
	paginationFunction func(*applicationautoscaling.DescribeScalableTargetsOutput, bool) bool,
) error {

	args := m.Called(in)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	paginationFunction(ExampleDescribeScalableTargetsOutput, true)
	return args.Error(0)
}
