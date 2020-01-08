package handlers

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
	"time"

	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
	complianceops "github.com/panther-labs/panther/api/gateway/compliance/client/operations"
	compliancemodels "github.com/panther-labs/panther/api/gateway/compliance/models"
)

const complianceCacheDuration = time.Minute

type complianceCacheEntry struct {
	ExpiresAt time.Time
	Policies  map[models.ID]*complianceStatus
}

type complianceStatus struct {
	SortIndex int                     // 0 is the top failing resource, 1 the next most failing, etc
	Status    models.ComplianceStatus // PASS, FAIL, ERROR, UNKNOWN
}

// Compliance-api queries are not very efficient right now.
// We cache the results here for a short period of time to minimize time the end-user is blocked
// when loading the policy detail page.
var complianceCache *complianceCacheEntry

// Get the pass/fail compliance status for a particular policy.
//
// Each org's pass/fail information for all policies is cached for a minute.
func getComplianceStatus(policyID models.ID) (*complianceStatus, error) {
	entry, err := getOrgCompliance()
	if err != nil {
		return nil, err
	}

	if result := entry.Policies[policyID]; result != nil {
		return result, nil
	}

	// A policy with no compliance entries is passing (it didn't evaluate anything)
	return &complianceStatus{SortIndex: -1, Status: models.ComplianceStatusPASS}, nil
}

func getOrgCompliance() (*complianceCacheEntry, error) {
	if complianceCache != nil && complianceCache.ExpiresAt.After(time.Now()) {
		return complianceCache, nil
	}

	zap.L().Info("loading policy pass/fail from compliance-api")
	result, err := complianceClient.Operations.DescribeOrg(&complianceops.DescribeOrgParams{
		Type:       "policy",
		HTTPClient: httpClient,
	})
	if err != nil {
		zap.L().Error("failed to load policy pass/fail from compliance-api", zap.Error(err))
		return nil, err
	}

	entry := &complianceCacheEntry{
		ExpiresAt: time.Now().Add(complianceCacheDuration),
		Policies:  make(map[models.ID]*complianceStatus, len(result.Payload.Policies)),
	}
	for i, policy := range result.Payload.Policies {
		entry.Policies[models.ID(*policy.ID)] = &complianceStatus{
			SortIndex: i,
			Status:    models.ComplianceStatus(policy.Status),
		}
	}
	complianceCache = entry
	return entry, nil
}

// Delete compliance status for entire policies or just some resource types within each policy.
func complianceBatchDelete(policies []*models.DeleteEntry, resourceTypes []string) error {
	entries := make([]*compliancemodels.DeleteStatus, len(policies))
	for i, policy := range policies {
		entries[i] = &compliancemodels.DeleteStatus{
			Policy: &compliancemodels.DeletePolicy{
				ID:            compliancemodels.PolicyID(policy.ID),
				ResourceTypes: resourceTypes,
			},
		}
	}

	zap.L().Info("deleting compliance entries", zap.Int("itemCount", len(entries)))
	_, err := complianceClient.Operations.DeleteStatus(&complianceops.DeleteStatusParams{
		Body:       &compliancemodels.DeleteStatusBatch{Entries: entries},
		HTTPClient: httpClient,
	})
	if err != nil {
		zap.L().Error("failed to delete compliance status", zap.Error(err))
		return err
	}

	// Remove the cached status as well
	for _, policy := range policies {
		if complianceCache != nil {
			delete(complianceCache.Policies, policy.ID)
		}
	}
	return nil
}

// Some policy changes trigger immediate or near-immediate changes to the compliance status.
//
// Examples:
//    - disabled policy => delete all associated compliance status
//    - changed python body => queue policy for full re-analysis
func updateComplianceStatus(oldItem, newItem *tableItem) error {
	if !newItem.Enabled {
		if oldItem != nil && oldItem.Enabled {
			zap.L().Info("policy is now disabled - deleting compliance status",
				zap.String("policyId", string(newItem.ID)))
			return complianceBatchDelete(
				[]*models.DeleteEntry{{ID: newItem.ID}}, []string{})
		}

		zap.L().Debug("policy remains disabled - no compliance updates required")
		return nil
	}

	// Delete compliance status for resource types which no longer apply
	if oldItem != nil && len(newItem.ResourceTypes) > 0 {
		if deletedTypes := setDifference(oldItem.ResourceTypes, newItem.ResourceTypes); len(deletedTypes) > 0 {
			zap.L().Info("policy no longer applies to some resource types - deleting compliance status",
				zap.String("policyId", string(newItem.ID)),
				zap.Strings("deletedTypes", deletedTypes))
			entries := []*models.DeleteEntry{{ID: newItem.ID}}
			if err := complianceBatchDelete(entries, deletedTypes); err != nil {
				return err
			}
		}
	}

	// Some changes require re-evaluating the entire policy
	if oldItem == nil || !oldItem.Enabled || oldItem.Body != newItem.Body || // newly enabled or updated Python
		(len(oldItem.ResourceTypes) > 0 && len(newItem.ResourceTypes) == 0) || // all resource types
		len(setDifference(newItem.ResourceTypes, oldItem.ResourceTypes)) > 0 { // additional resource types

		return queuePolicy(newItem)
	}

	// At this point, we know the compliance value (PASS/FAIL) won't change for any (policy, resource) pairs.
	// In other words, we don't need to re-evaluate the policy with the Python engine.
	//
	// But the compliance table has columns for severity and suppression -
	// if either of those changed, we can update the compliance API directly.
	if oldItem.Severity != newItem.Severity || !setEquality(oldItem.Suppressions, newItem.Suppressions) {
		return updateComplianceMetadata(newItem)
	}

	zap.L().Debug("policy has no major changes - no compliance updates required")
	return nil
}

// Update compliance status entries directly.
//
// This is used when only the policy severity / suppressions change - we don't need to rescan
// all affected resources in this case.
func updateComplianceMetadata(policy *tableItem) error {
	zap.L().Info("updating compliance status entry",
		zap.String("policyId", string(policy.ID)),
	)
	_, err := complianceClient.Operations.UpdateMetadata(&complianceops.UpdateMetadataParams{
		Body: &compliancemodels.UpdateMetadata{
			PolicyID:     compliancemodels.PolicyID(policy.ID),
			Severity:     compliancemodels.PolicySeverity(policy.Severity),
			Suppressions: compliancemodels.IgnoreSet(policy.Suppressions),
		},
		HTTPClient: httpClient,
	})
	return err
}
