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
