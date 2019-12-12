package api

import models "github.com/panther-labs/panther/api/snapshot"

// ListIntegrations returns all enabled integrations across each organization.
//
// The output of this handler is used to schedule pollers.
func (API) ListIntegrations(
	input *models.ListIntegrationsInput) ([]*models.SourceIntegration, error) {

	return db.ScanEnabledIntegrations(input)
}
