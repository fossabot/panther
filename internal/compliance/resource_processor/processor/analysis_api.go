package processor

import (
	"time"

	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/analysis/client/operations"
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
