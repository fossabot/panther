package api

import (
	"encoding/hex"

	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
)

// GetEvent retrieves a specific event
func (API) GetEvent(input *models.GetEventInput) (*models.GetEventOutput, error) {
	zap.L().Info("getting alert", zap.Any("input", input))

	binaryEventID, err := hex.DecodeString(*input.EventID)
	if err != nil {
		return nil, err
	}
	event, err := alertsDB.GetEvent(binaryEventID)
	if err != nil {
		return nil, err
	}

	return &models.GetEventOutput{
		Event: event,
	}, nil
}
