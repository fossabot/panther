package api

import (
	"github.com/panther-labs/panther/api/lambda/snapshot/models"
)

// DeleteIntegration deletes a specific integration.
func (API) DeleteIntegration(input *models.DeleteIntegrationInput) error {
	return db.DeleteIntegrationItem(input)
}
