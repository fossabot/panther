package poller

// ScanMsg contains a list of Scan Entries.
type ScanMsg struct {
	Entries []*ScanEntry `json:"entries"`
}

// ScanEntry indicates what type of scan should be performed, and provides the information needed
// to carry out that scan.
// The poller can scan a single resource, all resources of a given type, or all resources.
// Scanning all resources in an account is discouraged for performance reasons.
type ScanEntry struct {
	AWSAccountID     *string `json:"awsAccountId"`
	IntegrationID    *string `json:"integrationId"`
	Region           *string `json:"region"`
	ResourceID       *string `json:"resourceId"`
	ResourceType     *string `json:"resourceType"`
	ScanAllResources *bool   `json:"scanAllResources"`
}
