package api

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/api/lambda/organization/models"
)

// CreateOrganization generates a new organization ID.
//
// TODO - populate the rules table for new customers
func (API) CreateOrganization(
	input *models.CreateOrganizationInput) (*models.CreateOrganizationOutput, error) {

	// Then write the new org to the Dynamo table
	org := &models.Organization{
		AlertReportFrequency: input.AlertReportFrequency,
		AwsConfig:            input.AwsConfig,
		CreatedAt:            aws.String(time.Now().Format(time.RFC3339)),
		DisplayName:          input.DisplayName,
		Email:                input.Email,
		Phone:                input.Phone,
		RemediationConfig:    input.RemediationConfig,
	}

	if err := orgTable.Put(org); err != nil {
		return nil, err
	}
	return &models.CreateOrganizationOutput{Organization: org}, nil
}
