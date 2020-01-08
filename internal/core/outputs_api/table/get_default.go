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
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// GetDefault gets the default outputs for a given severity
func (table *DefaultsTable) GetDefault(severity *string) (*models.DefaultOutputsItem, error) {
	result, err := table.client.GetItem(
		&dynamodb.GetItemInput{
			TableName: table.Name,
			Key: DynamoItem{
				"severity": {S: severity},
			},
		})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, err
	}
	var defaultOutput models.DefaultOutputsItem
	if err = dynamodbattribute.UnmarshalMap(result.Item, &defaultOutput); err != nil {
		return nil, err
	}
	return &defaultOutput, nil
}
