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
