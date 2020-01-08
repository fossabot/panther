package handlers

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
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/compliance/models"
	"github.com/panther-labs/panther/pkg/awsbatch/dynamodbbatch"
)

// DeleteStatus deletes a batch of items
func DeleteStatus(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseDeleteStatus(request)
	if err != nil {
		return badRequest(err)
	}

	var deleteRequests []*dynamodb.WriteRequest
	for _, entry := range input.Entries {
		var entryRequests []*dynamodb.WriteRequest
		if entry.Policy != nil {
			entryRequests, err = policyDeleteEntries(entry.Policy.ID, entry.Policy.ResourceTypes)
		} else {
			entryRequests, err = resourceDeleteEntries(entry.Resource.ID)
		}

		if err != nil {
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}
		deleteRequests = append(deleteRequests, entryRequests...)
	}

	if len(deleteRequests) == 0 {
		// nothing to do
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
	}

	batchInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{Env.ComplianceTable: deleteRequests},
	}

	zap.L().Info("deleting batch of items", zap.Int("itemCount", len(deleteRequests)))
	if err := dynamodbbatch.BatchWriteItem(dynamoClient, maxWriteBackoff, batchInput); err != nil {
		zap.L().Error("dynamodbbatch.BatchWriteItem failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

func parseDeleteStatus(request *events.APIGatewayProxyRequest) (*models.DeleteStatusBatch, error) {
	var result models.DeleteStatusBatch
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	for i, entry := range result.Entries {
		if (entry.Resource == nil && entry.Policy == nil) || (entry.Resource != nil && entry.Policy != nil) {
			return nil, fmt.Errorf("entries[%d] invalid: exactly one of resource or policy is required", i)
		}
	}

	return &result, result.Validate(nil)
}

// Query the table for entries with the given policyID and return the list of delete requests.
func policyDeleteEntries(policyID models.PolicyID, resourceTypes []string) ([]*dynamodb.WriteRequest, error) {
	zap.L().Info("querying for deletion",
		zap.String("policyId", string(policyID)))
	keyCondition := expression.Key("policyId").Equal(expression.Value(policyID))
	projection := expression.NamesList(expression.Name("resourceId"))
	builder := expression.NewBuilder().WithKeyCondition(keyCondition).WithProjection(projection)

	// Filter the entries to just those of a specific resource type
	if len(resourceTypes) > 0 {
		var filter expression.ConditionBuilder

		for i, resourceType := range resourceTypes {
			typeFilter := expression.Equal(expression.Name("resourceType"), expression.Value(resourceType))
			if i == 0 {
				filter = typeFilter
			} else {
				filter = filter.Or(typeFilter)
			}
		}

		builder = builder.WithFilter(filter)
	}

	expr, err := builder.Build()
	if err != nil {
		zap.L().Error("expression.Build failed", zap.Error(err))
		return nil, err
	}

	// NOTE: You can't do a consistent read on a global index
	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		IndexName:                 &Env.IndexName,
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &Env.ComplianceTable,
	}

	var deleteRequests []*dynamodb.WriteRequest
	err = queryPages(input, func(item *models.ComplianceStatus) error {
		deleteRequests = append(deleteRequests, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{Key: tableKey(item.ResourceID, policyID)},
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return deleteRequests, nil
}

// Query the table for entries with the given resourceID and return the list of delete requests.
func resourceDeleteEntries(resourceID models.ResourceID) ([]*dynamodb.WriteRequest, error) {
	zap.L().Info("querying for deletion",
		zap.String("resourceId", string(resourceID)))
	keyCondition := expression.Key("resourceId").Equal(expression.Value(resourceID))
	projection := expression.NamesList(expression.Name("policyId"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).WithProjection(projection).Build()
	if err != nil {
		zap.L().Error("expression.Build failed", zap.Error(err))
		return nil, err
	}

	input := &dynamodb.QueryInput{
		ConsistentRead:            aws.Bool(true),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &Env.ComplianceTable,
	}

	var deleteRequests []*dynamodb.WriteRequest
	err = queryPages(input, func(item *models.ComplianceStatus) error {
		deleteRequests = append(deleteRequests, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{Key: tableKey(resourceID, item.PolicyID)},
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return deleteRequests, nil
}
