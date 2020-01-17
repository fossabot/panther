package table

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ListAlerts returns (a page of alerts, last evaluated key, any error)
func (table *AlertsTable) ListAlerts(exclusiveStartKey *string, pageSize *int) (
	summaries []*models.AlertItem, lastEvaluatedKey *string, err error) {

	var scanLimit *int64
	if pageSize != nil {
		scanLimit = aws.Int64(int64(*pageSize))
	}

	var scanExclusiveStartKey map[string]*dynamodb.AttributeValue
	if exclusiveStartKey != nil {
		paginationKey := &listAlertsPaginationKey{}
		err = jsoniter.UnmarshalFromString(*exclusiveStartKey, paginationKey)
		if err != nil {
			return nil, nil, err
		}
		scanExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"creationTime": {S: paginationKey.CreationTime},
			"alertId":      {S: paginationKey.AlertID},
		}
	}

	var scanInput = &dynamodb.ScanInput{
		TableName:         aws.String(table.AlertsTableName),
		ExclusiveStartKey: scanExclusiveStartKey,
		Limit:             scanLimit,
	}

	// TODO: Sort this by time (scan does not guarantee sortedness)
	scanOutput, err := table.Client.Scan(scanInput)
	if err != nil {
		return nil, nil, &genericapi.InternalError{
			Message: "scan to DDB failed" + err.Error(),
		}
	}

	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &summaries)
	if err != nil {
		return nil, nil, &genericapi.InternalError{
			Message: "failed to marshall response" + err.Error(),
		}
	}

	// If DDB returned a LastEvaluatedKey, it means there are more alerts to be returned
	// Return populated `lastEvaluatedKey` in the response.
	if len(scanOutput.LastEvaluatedKey) > 0 {
		paginationKey := listAlertsPaginationKey{
			CreationTime: scanOutput.LastEvaluatedKey["creationTime"].S,
			AlertID:      scanOutput.LastEvaluatedKey["alertId"].S,
		}
		marshalledKey, err := jsoniter.MarshalToString(paginationKey)
		if err != nil {
			return nil, nil, &genericapi.InternalError{
				Message: "failed to marshall key" + err.Error(),
			}
		}
		lastEvaluatedKey = &marshalledKey
	}

	return summaries, lastEvaluatedKey, nil
}

type listAlertsPaginationKey struct {
	CreationTime *string `json:"creationTime"`
	AlertID      *string `json:"alertId"`
}
