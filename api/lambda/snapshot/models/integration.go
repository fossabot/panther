package models

import "time"

// SourceIntegration is the dynamodb item corresponding to the PutIntegration route.
type SourceIntegration struct {
	*SourceIntegrationMetadata
	*SourceIntegrationStatus
	*SourceIntegrationScanInformation
}

// SourceIntegrationMetadata is general settings and metadata for an integration.
type SourceIntegrationMetadata struct {
	AWSAccountID         *string    `json:"awsAccountId"`
	CreatedAtTime        *time.Time `json:"createdAtTime"`
	CreatedBy            *string    `json:"createdBy"`
	IntegrationID        *string    `json:"integrationId"`
	IntegrationLabel     *string    `json:"integrationLabel"`
	IntegrationType      *string    `json:"integrationType"`
	ScanEnabled          *bool      `json:"scanEnabled"`
	ScanIntervalMins     *int       `json:"scanIntervalMins"`
	SourceSnsTopicArn    *string    `json:"sourceSnsTopicArn"`
	LogProcessingRoleArn *string    `json:"logProcessingRoleArn"`
}

// SourceIntegrationStatus provides context that the full scan works and that events are being received.
type SourceIntegrationStatus struct {
	ScanStatus  *string `json:"scanStatus"`
	EventStatus *string `json:"eventStatus"`
}

// SourceIntegrationScanInformation is detail about the last snapshot.
type SourceIntegrationScanInformation struct {
	LastScanEndTime      *time.Time `json:"lastScanEndTime"`
	LastScanErrorMessage *string    `json:"lastScanErrorMessage"`
	LastScanStartTime    *time.Time `json:"lastScanStartTime"`
}
