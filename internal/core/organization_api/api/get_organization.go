package api

import "github.com/panther-labs/panther/api/lambda/organization/models"

// GetOrganization retrieves customer account details.
func (API) GetOrganization(_ *models.GetOrganizationInput) (*models.GetOrganizationOutput, error) {
	org, err := orgTable.Get()
	if err != nil {
		return nil, err
	}

	return &models.GetOrganizationOutput{Organization: org}, nil
}
