package gateway

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"

	"github.com/panther-labs/panther/api/lambda/onboarding/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// DescribeExecution calls step function api to get execution status
func (g *StepFunctionGateway) DescribeExecution(executionArn *string) (*models.GetOnboardingStatusOutput, error) {
	eo, err := g.sfnClient.DescribeExecution(&sfn.DescribeExecutionInput{
		ExecutionArn: executionArn,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "sfn.describeExecution", Err: err}
	}
	return &models.GetOnboardingStatusOutput{
		Status:    eo.Status,
		StartDate: aws.String(eo.StartDate.String()),
		StopDate:  aws.String(eo.StopDate.String()),
	}, nil
}
