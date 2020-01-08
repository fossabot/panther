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

const (
	// IntegrationTypeAWSScan is the integration type for snapshots in customer AWS accounts.
	IntegrationTypeAWSScan = "aws-scan"
	// IntegrationTypeAWS3 is the integration type for importing data from customer S3 buckets.
	IntegrationTypeAWS3 = "aws-s3"

	// StatusError is the string set in the database when an error occurs in a scan.
	StatusError = "error"
	// StatusOK is the string set in the database when a scan is successful.
	StatusOK = "ok"
	// StatusScanning is the status set while a scan is underway.
	StatusScanning = "scanning"
)
