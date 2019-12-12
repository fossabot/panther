package api

import (
	models "github.com/panther-labs/panther/api/snapshot"
)

// DeleteIntegration deletes a specific integration.
func (API) DeleteIntegration(input *models.DeleteIntegrationInput) error {
	return db.DeleteIntegrationItem(input)
}
