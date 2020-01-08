package poller

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
