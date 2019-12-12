package api

import (
	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// SetDefaultOutputs sets the default outputs for an organization
func (API) SetDefaultOutputs(input *models.SetDefaultOutputsInput) (output *models.SetDefaultOutputsOutput, err error) {
	// Verify that the outputsIds exist
	for _, outputID := range input.OutputIDs {
		if _, err = outputsTable.GetOutput(outputID); err != nil {
			return nil, err
		}
	}

	item := &models.DefaultOutputsItem{
		Severity:  input.Severity,
		OutputIDs: input.OutputIDs,
	}

	if err = defaultsTable.PutDefaults(item); err != nil {
		return nil, err
	}

	output = &models.DefaultOutputs{
		Severity:  input.Severity,
		OutputIDs: input.OutputIDs,
	}

	return output, err
}
