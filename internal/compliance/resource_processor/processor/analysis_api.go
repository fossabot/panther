package processor

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

	"github.com/panther-labs/panther/api/gateway/analysis/client/operations"
)

const cacheDuration = 30 * time.Second

type policyCacheEntry struct {
	LastUpdated time.Time
	Policies    policyMap
}

var policyCache policyCacheEntry

// Get enabled policies from either the memory cache or the policy-api
func getPolicies() (policyMap, error) {
	if policyCache.Policies != nil && policyCache.LastUpdated.Add(cacheDuration).After(time.Now()) {
		// Cache entry exists and hasn't expired yet
		zap.L().Info("using policy cache",
			zap.Int("policyCount", len(policyCache.Policies)))
		return policyCache.Policies, nil
	}

	// Load from policy-api
	result, err := analysisClient.Operations.GetEnabledPolicies(
		&operations.GetEnabledPoliciesParams{HTTPClient: httpClient})
	if err != nil {
		zap.L().Error("failed to load policies from policy-api", zap.Error(err))
		return nil, err
	}
	zap.L().Info("successfully loaded enabled policies from policy-api",
		zap.Int("policyCount", len(result.Payload.Policies)))

	// Convert list of policies into a map by ID
	policies := make(policyMap, len(result.Payload.Policies))
	for _, policy := range result.Payload.Policies {
		policies[string(policy.ID)] = policy
	}

	policyCache = policyCacheEntry{LastUpdated: time.Now(), Policies: policies}
	return policies, nil
}
