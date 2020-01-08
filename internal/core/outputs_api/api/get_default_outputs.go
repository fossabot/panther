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
