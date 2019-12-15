package models

// Action defines an action the organization took
type Action = string

const (
	// VisitedOnboardingFlow defines when an organization visited the onboarding flow
	VisitedOnboardingFlow Action = "VISITED_ONBOARDING_FLOW"
)

// Organization defines the fields in the table row.
type Organization struct {
	AlertReportFrequency *string            `json:"alertReportFrequency"`
	AwsConfig            *AwsConfig         `json:"awsConfig"`
	CompletedActions     []*Action          `dynamodbav:"completedActions,omitempty,stringset" json:"completedActions"`
	CreatedAt            *string            `json:"createdAt"`
	DisplayName          *string            `json:"displayName"`
	Email                *string            `json:"email"`
	Phone                *string            `json:"phone"`
	RemediationConfig    *RemediationConfig `json:"remediationConfig,omitempty"`
}

// AwsConfig defines metadata related to AWS infrastructure for the organization
type AwsConfig struct {
	UserPoolID     *string `json:"userPoolId"`
	AppClientID    *string `json:"appClientId"`
	IdentityPoolID *string `json:"identityPoolId"`
}

// RemediationConfig contains information related to Remediation actions
type RemediationConfig struct {
	// Each organization will have one Lambda that is able to perform remediation for their AWS infrastructure.
	// This field contains the ARN for that Lambda.
	AwsRemediationLambdaArn *string `json:"awsRemediationLambdaArn,omitempty"`
}
