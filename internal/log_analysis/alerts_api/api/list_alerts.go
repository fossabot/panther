package api

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
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/analysis/client/operations"
	"github.com/panther-labs/panther/api/lambda/alerts/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// ListAlerts retrieves alert and event details.
func (API) ListAlerts(input *models.ListAlertsInput) (result *models.ListAlertsOutput, err error) {
	zap.L().Info("listing alerts", zap.Any("input", input))

	result = &models.ListAlertsOutput{}
	var alertItems []*models.AlertItem
	if input.RuleID != nil {
		zap.L().Info("fetching alert summaries for rule",
			zap.String("ruleId", *input.RuleID))
		alertItems, result.LastEvaluatedKey, err = alertsDB.ListAlertsByRule(input.RuleID, input.ExclusiveStartKey, input.PageSize)
	} else {
		zap.L().Info("fetching all alert summaries")
		alertItems, result.LastEvaluatedKey, err = alertsDB.ListAlerts(input.ExclusiveStartKey, input.PageSize)
	}
	if err != nil {
		return nil, err
	}

	result.Alerts, err = alertItemsToAlertSummary(alertItems)
	if err != nil {
		return nil, err
	}

	gatewayapi.ReplaceMapSliceNils(result)
	return result, nil
}

// alertItemsToAlertSummary converts a DDB Alert Item to an Alert Summary that will be returned by the API
func alertItemsToAlertSummary(items []*models.AlertItem) ([]*models.AlertSummary, error) {
	result := make([]*models.AlertSummary, len(items))

	// Many of the alerts returned might be triggered from the same rule
	// We are going to use this map in order to get the unique ruleIds
	ruleIDToSeverity := make(map[string]*string)

	for i, item := range items {
		ruleIDToSeverity[*item.RuleID] = nil
		result[i] = &models.AlertSummary{
			AlertID:          item.AlertID,
			RuleID:           item.RuleID,
			CreationTime:     item.CreationTime,
			LastEventMatched: item.LastEventMatched,
			EventsMatched:    aws.Int(len(item.EventHashes)),
		}
	}

	// Get the severity of each rule ID
	for ruleID := range ruleIDToSeverity {
		// All items are for the same org
		severity, err := getSeverity(aws.String(ruleID))
		if err != nil {
			return nil, err
		}
		ruleIDToSeverity[ruleID] = severity
	}

	// Set the correct severity
	for _, summary := range result {
		summary.Severity = ruleIDToSeverity[*summary.RuleID]
	}
	return result, nil
}

// getSeverity retrieves the rule severity associated with an alert
func getSeverity(ruleID *string) (*string, error) {
	zap.L().Debug("fetching severity of rule",
		zap.String("ruleId", *ruleID))

	response, err := policiesClient.Operations.GetRule(&operations.GetRuleParams{
		RuleID:     *ruleID,
		HTTPClient: httpClient,
	})

	if err != nil {
		return nil, err
	}

	return aws.String(string(response.Payload.Severity)), nil
}
