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
