package api

import "github.com/panther-labs/panther/api/lambda/organization/models"

// UpdateOrganization updates account details.
func (API) UpdateOrganization(
	input *models.UpdateOrganizationInput) (*models.UpdateOrganizationOutput, error) {

	updated, err := orgTable.Update(&models.Organization{
		AlertReportFrequency: input.AlertReportFrequency,
		AwsConfig:            input.AwsConfig,
		DisplayName:          input.DisplayName,
		Email:                input.Email,
		Phone:                input.Phone,
		RemediationConfig:    input.RemediationConfig,
	})
	if err != nil {
		return nil, err
	}

	return &models.UpdateOrganizationOutput{Organization: updated}, nil
}
