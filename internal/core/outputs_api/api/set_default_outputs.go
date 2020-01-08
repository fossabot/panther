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
