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
	"encoding/hex"

	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
)

// GetEvent retrieves a specific event
func (API) GetEvent(input *models.GetEventInput) (*models.GetEventOutput, error) {
	zap.L().Info("getting alert", zap.Any("input", input))

	binaryEventID, err := hex.DecodeString(*input.EventID)
	if err != nil {
		return nil, err
	}
	event, err := alertsDB.GetEvent(binaryEventID)
	if err != nil {
		return nil, err
	}

	return &models.GetEventOutput{
		Event: event,
	}, nil
}
