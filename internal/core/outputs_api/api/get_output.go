package api

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
