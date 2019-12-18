package api

import (
	"bytes"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// GetAlert retrieves details for a given alert
func (API) GetAlert(input *models.GetAlertInput) (result *models.GetAlertOutput, err error) {
	zap.L().Info("getting alert", zap.Any("input", input))

	alertItem, err := alertsDB.GetAlert(input.AlertID)
	if err != nil {
		return nil, err
	}

	result = &models.Alert{
		AlertID:          alertItem.AlertID,
		RuleID:           alertItem.RuleID,
		CreationTime:     alertItem.CreationTime,
		LastEventMatched: alertItem.LastEventMatched,
		MatchedEventNum:  aws.Int(len(alertItem.EventHashes)),
	}

	var eventHashesToReturn [][]byte

	if input.EventsPageSize == nil { // if no eventsPageSize is defined, fallback to returning every event.
		eventHashesToReturn = alertItem.EventHashes
	} else {
		if input.EventsExclusiveStartKey == nil {
			// if no EventsExclusiveStartKey is defined, return the first events
			endIndex := len(alertItem.EventHashes)

			// If the full events are more than the page size
			// just return as many events as the page size
			// and set the lastEvaluatedKey
			if endIndex > *input.EventsPageSize {
				endIndex = *input.EventsPageSize
				lastHashToReturn := hex.EncodeToString(alertItem.EventHashes[endIndex-1])
				result.EventsLastEvaluatedKey = &lastHashToReturn
			}
			eventHashesToReturn = alertItem.EventHashes[:endIndex]
		} else {
			exclusiveStartKey, err := hex.DecodeString(*input.EventsExclusiveStartKey)
			if err != nil {
				return nil, err
			}
			for i, hash := range alertItem.EventHashes {
				if bytes.Equal(exclusiveStartKey, hash) {
					// We add 1 because the key is exclusive
					startIndex := i + 1
					endIndex := len(alertItem.EventHashes)

					// If the full events are more than the page size
					// just return as many events as the page size
					// and set the lastEvaluatedKey
					if endIndex > startIndex+*input.EventsPageSize {
						endIndex = startIndex + *input.EventsPageSize
						lastHashToReturn := hex.EncodeToString(alertItem.EventHashes[endIndex-1])
						result.EventsLastEvaluatedKey = &lastHashToReturn
					}
					eventHashesToReturn = alertItem.EventHashes[startIndex:endIndex]
					break
				}
			}
		}
	}

	for _, hash := range eventHashesToReturn {
		newEvent, err := alertsDB.GetEvent(hash)
		if err != nil {
			return nil, err
		}
		result.Events = append(result.Events, newEvent)
	}

	gatewayapi.ReplaceMapSliceNils(result)
	return result, nil
}
