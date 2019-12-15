package api

import (
	"github.com/panther-labs/panther/api/lambda/onboarding/models"
)

// GetOnboardingStatus calls stepFunctionGateway to get onboarding status.
func (API) GetOnboardingStatus(input *models.GetOnboardingStatusInput) (*models.GetOnboardingStatusOutput, error) {
	o, err := stepFunctionGateway.DescribeExecution(input.ExecutionArn)
	if err != nil {
		return nil, err
	}
	return o, nil
}
