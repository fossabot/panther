package models

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
