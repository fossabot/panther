package ddb

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
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/awsbatch/dynamodbbatch"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var maxElapsedTime = 15 * time.Second

// BatchPutSourceIntegrations adds a batch of new Snapshot Integrations to the database.
func (ddb *DDB) BatchPutSourceIntegrations(input []*models.SourceIntegrationMetadata) error {
	writeRequests := make([]*dynamodb.WriteRequest, len(input))

	// Marshal each new integration and add to the write request
	for i, item := range input {
		item, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return &genericapi.AWSError{Err: err, Method: "Dynamodb.MarshalMap"}
		}
		writeRequests[i] = &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: item}}
	}

	// Do the batch write
	err := dynamodbbatch.BatchWriteItem(ddb.Client, maxElapsedTime, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{ddb.TableName: writeRequests}})
	if err != nil {
		return &genericapi.AWSError{Err: err, Method: "Dynamodb.BatchWriteItem"}
	}

	return nil
}
