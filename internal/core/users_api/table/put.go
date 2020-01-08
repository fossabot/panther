package users

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
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// Put writes user to organization mapping to the dynamo table.
func (table *Table) Put(user *models.UserItem) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return &genericapi.InternalError{
			Message: "failed to marshal User to a dynamo item: " + err.Error()}
	}

	zap.L().Info("saving user to dynamo", zap.String("user email", *user.ID))
	input := &dynamodb.PutItemInput{Item: item, TableName: table.Name}
	if _, err = table.client.PutItem(input); err != nil {
		return &genericapi.AWSError{Method: "dynamodb.PutItem", Err: err}
	}

	return nil
}
