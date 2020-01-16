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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// GetIntegration returns an integration by its ID
func (ddb *DDB) GetIntegration(integrationID *string) (*models.SourceIntegrationMetadata, error) {
	output, err := ddb.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ddb.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			hashKey: {S: integrationID},
		},
	})
	if err != nil {
		return nil, &genericapi.AWSError{Err: err, Method: "Dynamodb.GetItem"}
	}

	var integration models.SourceIntegrationMetadata
	if output.Item == nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalMap(output.Item, &integration); err != nil {
		return nil, err
	}

	return &integration, nil
}
