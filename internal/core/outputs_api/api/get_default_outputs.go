package api

import (
	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// GetDefaultOutputs retrieves the default outputs for an organization
func (API) GetDefaultOutputs(input *models.GetDefaultOutputsInput) (result *models.GetDefaultOutputsOutput, err error) {
	items, err := defaultsTable.GetDefaults()
	if err != nil {
		return nil, err
	}

	defaults := []*models.DefaultOutputs{}
	for _, item := range items {
		if item.OutputIDs == nil {
			continue
		}
		outputs := &models.DefaultOutputs{
			Severity:  item.Severity,
			OutputIDs: item.OutputIDs,
		}
		defaults = append(defaults, outputs)
	}

	result = &models.GetDefaultOutputsOutput{Defaults: defaults}

	return result, nil
}
