package table

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
