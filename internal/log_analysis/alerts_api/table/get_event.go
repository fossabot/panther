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
)

// GetEvent retrieves an event from DDB
func (table *AlertsTable) GetEvent(eventHash []byte) (*string, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"eventHash": {B: eventHash},
		},
		TableName: aws.String(table.EventsTableName),
	}

	ddbResult, err := table.Client.GetItem(input)
	if err != nil {
		return nil, err
	}

	return ddbResult.Item["event"].S, nil
}
