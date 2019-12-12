package api

import (
	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// DeleteOutput removes the alert output configuration
func (API) DeleteOutput(input *models.DeleteOutputInput) error {
	defaults, err := defaultsTable.GetDefaults()
	if err != nil {
		return err
	}

	for _, defaultOutput := range defaults {
		for index, outputID := range defaultOutput.OutputIDs {
			if *outputID == *input.OutputID {
				if aws.BoolValue(input.Force) {
					// Remove outputID from the list of outputs
					ids := defaultOutput.OutputIDs
					defaultOutput.OutputIDs = append(ids[:index], ids[index+1:]...)

					// Update defaults table
					if err = defaultsTable.PutDefaults(defaultOutput); err != nil {
						return err
					}
				} else {
					return &genericapi.InUseError{Message: "This destination is currently in use, please try again in a few seconds"}
				}
			}
		}
	}

	return outputsTable.DeleteOutput(input.OutputID)
}
