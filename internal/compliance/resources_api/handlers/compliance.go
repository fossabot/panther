package handlers

import (
	"time"

	"go.uber.org/zap"

	complianceops "github.com/panther-labs/panther/api/compliance/client/operations"
	"github.com/panther-labs/panther/api/resources/models"
)

const complianceCacheDuration = time.Minute

type complianceCacheEntry struct {
	ExpiresAt time.Time
	Resources map[models.ResourceID]*complianceStatus
}

type complianceStatus struct {
	SortIndex int                     // 0 is the top failing resource, 1 the next most failing, etc
	Status    models.ComplianceStatus // PASS, FAIL, ERROR
}

var complianceCache *complianceCacheEntry

// Get the pass/fail compliance status for a particular resource.
//
// Each org's pass/fail information for all policies is cached for a minute.
func getComplianceStatus(resourceID models.ResourceID) (*complianceStatus, error) {
	entry, err := getOrgCompliance()
	if err != nil {
		return nil, err
	}

	if result := entry.Resources[resourceID]; result != nil {
		return result, nil
	}

	// A resource with no compliance entries is passing (no policies applied to it)
	return &complianceStatus{SortIndex: -1, Status: models.ComplianceStatusPASS}, nil
}

func getOrgCompliance() (*complianceCacheEntry, error) {
	if complianceCache != nil && complianceCache.ExpiresAt.After(time.Now()) {
		return complianceCache, nil
	}

	zap.L().Info("loading resource pass/fail from compliance-api")
	result, err := complianceClient.Operations.DescribeOrg(&complianceops.DescribeOrgParams{
		Type:       "resource",
		HTTPClient: httpClient,
	})
	if err != nil {
		zap.L().Error("failed to load resource pass/fail from compliance-api", zap.Error(err))
		return nil, err
	}

	entry := &complianceCacheEntry{
		ExpiresAt: time.Now().Add(complianceCacheDuration),
		Resources: make(map[models.ResourceID]*complianceStatus, len(result.Payload.Resources)),
	}
	for i, resource := range result.Payload.Resources {
		entry.Resources[models.ResourceID(*resource.ID)] = &complianceStatus{
			SortIndex: i,
			Status:    models.ComplianceStatus(resource.Status),
		}
	}
	complianceCache = entry
	return entry, nil
}
