package api

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	"github.com/panther-labs/panther/api/lambda/onboarding/models"
	"github.com/panther-labs/panther/internal/core/organization_onboarding/gateway"
	"github.com/panther-labs/panther/pkg/genericapi"
)

type mockGatewayStepFunctionClient struct {
	gateway.API
	stepFunctionGatewayErr bool
}

func (m *mockGatewayStepFunctionClient) DescribeExecution(executionArn *string) (*models.GetOnboardingStatusOutput, error) {
	if m.stepFunctionGatewayErr {
		return nil, &genericapi.AWSError{}
	}
	startDate, _ := time.Parse(time.RFC3339, "2019-04-10T23:00:00Z")
	stopDate, _ := time.Parse(time.RFC3339, "2019-04-10T22:59:00Zs")

	return &models.GetOnboardingStatusOutput{
		Status:    aws.String("PASSING"),
		StartDate: aws.String(startDate.String()),
		StopDate:  aws.String(stopDate.String()),
	}, nil
}

func TestGetOnboardingStatusGateway(t *testing.T) {
	stepFunctionGateway = &mockGatewayStepFunctionClient{}
	result, err := (API{}).GetOnboardingStatus(&models.GetOnboardingStatusInput{
		ExecutionArn: aws.String("fakeExecutionArns"),
	})
	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestGetOnboardingStatusGatewayErr(t *testing.T) {
	stepFunctionGateway = &mockGatewayStepFunctionClient{stepFunctionGatewayErr: true}
	result, err := (API{}).GetOnboardingStatus(&models.GetOnboardingStatusInput{
		ExecutionArn: aws.String("fakeExecutionArns"),
	})
	assert.Nil(t, result)
	assert.Error(t, err)
}
