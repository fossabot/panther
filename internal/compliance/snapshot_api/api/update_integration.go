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
	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb"
)

// UpdateIntegrationSettings makes an update to an integration from the UI.
//
// This endpoint updates attributes such as the behavior of the integration, or display information.
func (API) UpdateIntegrationSettings(input *models.UpdateIntegrationSettingsInput) error {
	return db.UpdateItem(&ddb.UpdateIntegrationItem{
		IntegrationID:    input.IntegrationID,
		IntegrationLabel: input.IntegrationLabel,
		ScanIntervalMins: input.ScanIntervalMins,
		ScanEnabled:      input.ScanEnabled,
		S3Buckets:        input.S3Buckets,
		KmsKeys:          input.KmsKeys,
	})
}

// UpdateIntegrationLastScanStart updates an integration when a new scan is started.
func (API) UpdateIntegrationLastScanStart(input *models.UpdateIntegrationLastScanStartInput) error {
	return db.UpdateItem(&ddb.UpdateIntegrationItem{
		IntegrationID:     input.IntegrationID,
		LastScanStartTime: input.LastScanStartTime,
		ScanStatus:        input.ScanStatus,
	})
}

// UpdateIntegrationLastScanEnd updates an integration when a scan ends.
func (API) UpdateIntegrationLastScanEnd(input *models.UpdateIntegrationLastScanEndInput) error {
	return db.UpdateItem(&ddb.UpdateIntegrationItem{
		IntegrationID:        input.IntegrationID,
		LastScanEndTime:      input.LastScanEndTime,
		LastScanErrorMessage: input.LastScanErrorMessage,
		ScanStatus:           input.ScanStatus,
	})
}
