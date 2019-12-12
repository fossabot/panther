package api

import (
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// GetOutput retrieves a single alert output
func (API) GetOutput(input *models.GetOutputInput) (*models.GetOutputOutput, error) {
	item, err := outputsTable.GetOutput(input.OutputID)
	if err != nil {
		return nil, err
	}

	defaults, err := defaultsTable.GetDefaults()
	if err != nil {
		return nil, err
	}

	alertOutput, err := populateAlertOutput(item, defaults)
	if err != nil {
		return nil, err
	}

	return alertOutput, nil
}

// Checks if an Alert Output is marked as verified or not.
// If the Alert Output is not marked as verified, we check the state of the
// configuration processes and update as appropriate
func checkAndUpdateVerificationStatus(output *models.AlertOutput) error {
	if *output.VerificationStatus == models.VerificationStatusSuccess {
		return nil
	}

	zap.L().Info("update the verification status of output",
		zap.String("outputId", *output.OutputID))
	newStatus, err := outputVerification.GetVerificationStatus(output)
	if err != nil {
		return err
	}
	if *newStatus != *output.VerificationStatus {
		zap.L().Info("verification status of output has changed",
			zap.String("oldVerificationStatus", *output.VerificationStatus),
			zap.String("newVerificationStatus", *newStatus))
		output.VerificationStatus = newStatus
		outputItem, err := AlertOutputToItem(output)
		if err != nil {
			return err
		}
		_, err = outputsTable.UpdateOutput(outputItem)
		if err != nil {
			return nil
		}
	}
	return nil
}

func populateAlertOutput(item *models.AlertOutputItem, defaultOutputs []*models.DefaultOutputsItem) (*models.AlertOutput, error) {
	alertOutput, err := ItemToAlertOutput(item)
	if err != nil {
		return nil, err
	}

	if err = checkAndUpdateVerificationStatus(alertOutput); err != nil {
		return nil, err
	}

	alertOutput.DefaultForSeverity = []*string{}
	for _, defaultOutput := range defaultOutputs {
		for _, outputID := range defaultOutput.OutputIDs {
			if *outputID == *alertOutput.OutputID {
				alertOutput.DefaultForSeverity = append(alertOutput.DefaultForSeverity, defaultOutput.Severity)
			}
		}
	}
	return alertOutput, nil
}
