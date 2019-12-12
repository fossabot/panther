package api

import (
	"github.com/panther-labs/panther/api/lambda/snapshot/models"
)

// UpdateIntegrationSettings makes an update to an integration from the UI.
//
// This endpoint updates attributes such as the behavior of the integration, or display information.
func (API) UpdateIntegrationSettings(input *models.UpdateIntegrationSettingsInput) error {
	return db.UpdateItem(&models.UpdateIntegrationItem{
		IntegrationID:        input.IntegrationID,
		IntegrationLabel:     input.IntegrationLabel,
		ScanIntervalMins:     input.ScanIntervalMins,
		ScanEnabled:          input.ScanEnabled,
		LogProcessingRoleArn: input.LogProcessingRoleArn,
		SourceSnsTopicArn:    input.SourceSnsTopicArn,
	})
}

// UpdateIntegrationLastScanStart updates an integration when a new scan is started.
func (API) UpdateIntegrationLastScanStart(input *models.UpdateIntegrationLastScanStartInput) error {
	return db.UpdateItem(&models.UpdateIntegrationItem{
		IntegrationID:     input.IntegrationID,
		LastScanStartTime: input.LastScanStartTime,
		ScanStatus:        input.ScanStatus,
	})
}

// UpdateIntegrationLastScanEnd updates an integration when a scan ends.
func (API) UpdateIntegrationLastScanEnd(input *models.UpdateIntegrationLastScanEndInput) error {
	return db.UpdateItem(&models.UpdateIntegrationItem{
		IntegrationID:        input.IntegrationID,
		LastScanEndTime:      input.LastScanEndTime,
		LastScanErrorMessage: input.LastScanErrorMessage,
		ScanStatus:           input.ScanStatus,
	})
}
