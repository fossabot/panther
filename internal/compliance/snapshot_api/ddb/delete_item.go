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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// DeleteIntegrationItem deletes an integration from the database based on the integration ID
func (ddb *DDB) DeleteIntegrationItem(input *models.DeleteIntegrationInput) error {
	condition := expression.AttributeExists(expression.Name("integrationId"))

	builder := expression.NewBuilder().WithCondition(condition)
	expr, err := builder.Build()
	if err != nil {
		return &genericapi.InternalError{Message: "failed to build DeleteIntegration ddb expression"}
	}

	_, err = ddb.Client.DeleteItem(&dynamodb.DeleteItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		Key: map[string]*dynamodb.AttributeValue{
			hashKey: {S: input.IntegrationID},
		},
		TableName: aws.String(ddb.TableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return &genericapi.DoesNotExistError{Message: aerr.Error()}
			default:
				return &genericapi.AWSError{Err: err, Method: "Dynamodb.DeleteItem"}
			}
		}
	}

	return nil
}
