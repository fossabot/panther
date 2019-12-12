package models

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
