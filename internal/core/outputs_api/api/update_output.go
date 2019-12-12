package api

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// UpdateOutput updates the alert output with the new values
func (API) UpdateOutput(input *models.UpdateOutputInput) (*models.UpdateOutputOutput, error) {
	existingOutput, err := outputsTable.GetOutputByName(input.DisplayName)
	if err != nil {
		return nil, err
	}

	// if there is already a different output with the same 'displayName', fail the operation
	if existingOutput != nil && *existingOutput.OutputID != *input.OutputID {
		return nil, &genericapi.AlreadyExistsError{
			Message: "A destination with the name" + *input.DisplayName + " already exists, please choose another display name"}
	}

	outputType, err := getOutputType(input.OutputConfig)
	if err != nil {
		return nil, err
	}

	alertOutput := &models.AlertOutput{
		DisplayName:        input.DisplayName,
		LastModifiedBy:     input.UserID,
		LastModifiedTime:   aws.String(time.Now().Format(time.RFC3339)),
		OutputID:           input.OutputID,
		OutputType:         outputType,
		OutputConfig:       input.OutputConfig,
		DefaultForSeverity: input.DefaultForSeverity,
	}

	alertOutputItem, err := AlertOutputToItem(alertOutput)
	if err != nil {
		return nil, err
	}

	if alertOutputItem, err = outputsTable.UpdateOutput(alertOutputItem); err != nil {
		return nil, err
	}

	defaults, err := defaultsTable.GetDefaults()
	if err != nil {
		return nil, err
	}

	// Removing outputId from all defaults
	for _, defaultOutput := range defaults {
		defaultOutput.OutputIDs = removeFromSlice(defaultOutput.OutputIDs, input.OutputID)
		if err := defaultsTable.PutDefaults(defaultOutput); err != nil {
			return nil, err
		}
	}

	if err := addToDefaults(input.DefaultForSeverity, input.OutputID); err != nil {
		return nil, err
	}

	alertOutput.CreatedBy = alertOutputItem.CreatedBy
	alertOutput.CreationTime = alertOutputItem.CreationTime
	alertOutput.VerificationStatus = alertOutputItem.VerificationStatus

	return alertOutput, nil
}

func removeFromSlice(slice []*string, item *string) []*string {
	for i, element := range slice {
		if *element == *item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
