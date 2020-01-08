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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// DeleteOutput removes an output from the table.
func (table *OutputsTable) DeleteOutput(outputID *string) error {
	condition := expression.Name("outputId").Equal(expression.Value(outputID))

	conditionExpression, err := expression.NewBuilder().WithCondition(condition).Build()

	if err != nil {
		return &genericapi.InternalError{Message: "failed to build expression " + err.Error()}
	}

	_, err = table.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: table.Name,
		Key: DynamoItem{
			"outputId": {S: outputID},
		},
		ConditionExpression:       conditionExpression.Condition(),
		ExpressionAttributeNames:  conditionExpression.Names(),
		ExpressionAttributeValues: conditionExpression.Values(),
	})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return &genericapi.DoesNotExistError{Message: "outputId=" + *outputID}
		}
		return &genericapi.AWSError{Method: "dynamodb.DeleteItem", Err: err}
	}

	return nil
}
