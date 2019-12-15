package models

// LambdaInput is the request structure for the organization-api Lambda function.
type LambdaInput struct {
	CompleteAction     *CompleteActionInput     `json:"getCompletedActions"`
	CreateOrganization *CreateOrganizationInput `json:"createOrganization"`
	GetOrganization    *GetOrganizationInput    `json:"getOrganization"`
	UpdateOrganization *UpdateOrganizationInput `json:"updateOrganization"`
}

// CompleteActionInput Adds a Action to an Organization
type CompleteActionInput struct {
	CompletedActions []*Action `json:"actions"`
}

// CompleteActionOutput Adds a Action to an Organization
type CompleteActionOutput struct {
	CompletedActions []*Action `json:"actions"`
}

// CreateOrganizationInput creates a new Panther customer account.
type CreateOrganizationInput struct {
	AlertReportFrequency *string            `json:"alertReportFrequency" validate:"omitempty,oneof=P1D P1W"`
	AwsConfig            *AwsConfig         `json:"awsConfig"`
	DisplayName          *string            `json:"displayName" validate:"required,min=1"`
	Email                *string            `genericapi:"redact" json:"email" validate:"required,email"`
	Phone                *string            `genericapi:"redact" json:"phone"`
	RemediationConfig    *RemediationConfig `json:"remediationConfig,omitempty"`
}

// CreateOrganizationOutput returns the newly created organization.
type CreateOrganizationOutput struct {
	Organization *Organization `json:"organization"`
}

// GetOrganizationInput retrieves the details of a Panther customer account.
type GetOrganizationInput struct {
}

// GetOrganizationOutput is the table row representing a customer account.
type GetOrganizationOutput struct {
	Organization *Organization `json:"organization"`
}

// UpdateOrganizationInput modifies the details of an existing organization.
type UpdateOrganizationInput struct {
	CreateOrganizationInput
}

// UpdateOrganizationOutput is the table row representing the modified customer account.
type UpdateOrganizationOutput struct {
	Organization *Organization `json:"organization"`
}
