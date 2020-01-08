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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// Update updates account details and returns the updated item
func (table *OrganizationsTable) Update(org *models.Organization) (*models.Organization, error) {
	update := expression.
		Set(expression.Name("alertReportFrequency"), expression.Value(org.AlertReportFrequency)).
		Set(expression.Name("awsConfig"), expression.Value(org.AwsConfig)).
		Set(expression.Name("displayName"), expression.Value(org.DisplayName)).
		Set(expression.Name("email"), expression.Value(org.Email)).
		Set(expression.Name("phone"), expression.Value(org.Phone)).
		Set(expression.Name("remediationConfig"), expression.Value(org.RemediationConfig))
	return table.doUpdate(update)
}

type flagSet []*models.Action

// Marshal string slice as a Dynamo StringSet instead of a List
func (s flagSet) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	av.SS = make([]*string, 0, len(s))
	for _, flag := range s {
		av.SS = append(av.SS, flag)
	}
	return nil
}

// AddActions append additional actions to completed actions and returns the updated organization
func (table *OrganizationsTable) AddActions(actions []*models.Action) (*models.Organization, error) {
	update := expression.Add(
		expression.Name("completedActions"), expression.Value(flagSet(actions)))
	return table.doUpdate(update)
}

func (table *OrganizationsTable) doUpdate(update expression.UpdateBuilder) (*models.Organization, error) {
	condition := expression.AttributeExists(expression.Name("id"))

	expr, err := expression.NewBuilder().WithCondition(condition).WithUpdate(update).Build()
	if err != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to build update expression: " + err.Error()}
	}

	input := &dynamodb.UpdateItemInput{
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       DynamoItem{"id": {S: aws.String("1")}},
		ReturnValues:              aws.String("ALL_NEW"),
		TableName:                 table.Name,
		UpdateExpression:          expr.Update(),
	}

	zap.L().Info("updating org in dynamo")
	response, err := table.client.UpdateItem(input)

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return nil, &genericapi.DoesNotExistError{}
		}
		return nil, &genericapi.AWSError{Method: "dynamodb.UpdateItem", Err: err}
	}

	var newOrg models.Organization
	if err = dynamodbattribute.UnmarshalMap(response.Attributes, &newOrg); err != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to unmarshal dynamo item to an Organization: " + err.Error()}
	}

	return &newOrg, nil
}
