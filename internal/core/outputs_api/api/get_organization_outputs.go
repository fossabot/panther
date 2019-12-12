package api

import (
	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// GetOrganizationOutputs returns all the alert outputs configured for one organization
func (API) GetOrganizationOutputs(input *models.GetOrganizationOutputsInput) (models.GetOrganizationOutputsOutput, error) {
	outputItems, err := outputsTable.GetOutputs()
	if err != nil {
		return nil, err
	}

	defaults, err := defaultsTable.GetDefaults()
	if err != nil {
		return nil, err
	}

	outputs := make([]*models.AlertOutput, len(outputItems))
	for i, item := range outputItems {
		alertOutput, err := populateAlertOutput(item, defaults)
		if err != nil {
			return nil, err
		}

		outputs[i] = alertOutput
	}

	return outputs, nil
}
