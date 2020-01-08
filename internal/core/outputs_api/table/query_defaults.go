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
	"github.com/panther-labs/panther/pkg/genericapi"
)

// GetDefaults returns the default outputs for one organization
func (table *DefaultsTable) GetDefaults() (defaults []*models.DefaultOutputsItem, err error) {
	var scanInput = &dynamodb.ScanInput{
		TableName: table.Name,
	}

	var internalErr error
	queryErr := table.client.ScanPages(scanInput,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			var defaultsPartial []*models.DefaultOutputsItem
			if internalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &defaultsPartial); internalErr != nil {
				internalErr = &genericapi.InternalError{
					Message: "failed to unmarshal dynamo item to an AlertOutputItem: " + internalErr.Error(),
				}
				return false
			}
			defaults = append(defaults, defaultsPartial...)
			return true
		})

	if queryErr != nil {
		return nil, &genericapi.AWSError{Err: queryErr, Method: "dynamodb.ScanPages"}
	}
	if internalErr != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to unmarshal dynamo items: " + internalErr.Error()}
	}

	return defaults, nil
}
