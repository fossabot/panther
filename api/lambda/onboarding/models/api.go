package models

// LambdaInput is the invocation event expected by the Lambda function.
//
// Exactly one action must be specified.
type LambdaInput struct {
	GetOnboardingStatus *GetOnboardingStatusInput `json:"getOnboardingStatus"`
}

// GetOnboardingStatusInput gets the step function status
type GetOnboardingStatusInput struct {
	ExecutionArn *string `json:"executionArn" validate:"required"`
}

// GetOnboardingStatusOutput returns the state machine status
type GetOnboardingStatusOutput struct {
	Status    *string `json:"status"`
	StartDate *string `json:"startDate"`
	StopDate  *string `json:"stopDate"`
}
