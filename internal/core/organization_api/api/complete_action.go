package api

import (
	"github.com/panther-labs/panther/api/lambda/organization/models"
)

// CompleteAction generates a new organization ID.
func (API) CompleteAction(input *models.CompleteActionInput) (*models.CompleteActionOutput, error) {
	updatedOrg, err := orgTable.AddActions(input.CompletedActions)
	if err != nil {
		return nil, err
	}
	return &models.CompleteActionOutput{CompletedActions: updatedOrg.CompletedActions}, nil
}
