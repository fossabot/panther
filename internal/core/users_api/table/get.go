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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// Get retrieves user to org mapping from the table.
func (table *Table) Get(id *string) (*models.UserItem, error) {
	zap.L().Info("retrieving user from dynamo", zap.String("id", *id))
	response, err := table.client.GetItem(&dynamodb.GetItemInput{
		Key:       DynamoItem{"id": {S: id}},
		TableName: table.Name,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "dynamodb.GetItem", Err: err}
	}

	var user models.UserItem
	if err = dynamodbattribute.UnmarshalMap(response.Item, &user); err != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to unmarshal dynamo item to a User: " + err.Error()}
	}

	if aws.StringValue(user.ID) == "" {
		return nil, &genericapi.DoesNotExistError{Message: "id=" + *id}
	}

	return &user, nil
}
