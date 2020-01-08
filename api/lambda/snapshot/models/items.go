package models

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

import "time"

// UpdateIntegrationItem updates almost every attribute in the table.
//
// It's used for attributes that can change, which is almost all of them except for the
// creation based ones (CreatedAtTime and CreatedBy).
type UpdateIntegrationItem struct {
	ScanEnabled          *bool      `json:"scanEnabled"`
	IntegrationID        *string    `json:"integrationId"`
	IntegrationLabel     *string    `json:"integrationLabel"`
	IntegrationType      *string    `json:"integrationType"`
	LastScanEndTime      *time.Time `json:"lastScanEndTime"`
	LastScanErrorMessage *string    `json:"lastScanErrorMessage"`
	LastScanStartTime    *time.Time `json:"lastScanStartTime"`
	ScanStatus           *string    `json:"scanStatus"`
	ScanIntervalMins     *int       `json:"scanIntervalMins"`
	SourceSnsTopicArn    *string    `json:"sourceSnsTopicArn"`
	LogProcessingRoleArn *string    `json:"logProcessingRoleArn"`
}
