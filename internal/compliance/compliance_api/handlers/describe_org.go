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
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/compliance/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// DescribeOrg returns pass/fail counts for every policy or resource in a customer account.
func DescribeOrg(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	queryInput, err := buildDescribeOrgScan()
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	var result models.EntireOrg
	if request.QueryStringParameters["type"] == "policy" {
		result.Policies, err = listPolicies(queryInput)
	} else {
		result.Resources, err = listResources(queryInput)
	}
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(&result, http.StatusOK)
}

func buildDescribeOrgScan() (*dynamodb.ScanInput, error) {
	filter := expression.Equal(expression.Name("suppressed"), expression.Value(false))
	projection := expression.NamesList(
		expression.Name("policyId"),
		expression.Name("policySeverity"),
		expression.Name("resourceId"),
		expression.Name("status"),
	)

	expr, err := expression.NewBuilder().
		WithFilter(filter).
		WithProjection(projection).
		Build()
	if err != nil {
		zap.L().Error("expression.Build failed", zap.Error(err))
		return nil, err
	}

	return &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &Env.ComplianceTable,
	}, nil
}

// List policies, sort by top failing
func listPolicies(input *dynamodb.ScanInput) ([]*models.ItemSummary, error) {
	// Fetch all policies from Dynamo
	policyMap, _, err := scanGroupByID(input, true, false)
	if err != nil {
		return nil, err
	}

	// Convert to a slice and sort by top failing
	policySlice := make([]*models.PolicySummary, 0, len(policyMap))
	for _, policy := range policyMap {
		policySlice = append(policySlice, policy)
	}
	sortPoliciesByTopFailing(policySlice)

	// Convert to final ItemSummary
	result := make([]*models.ItemSummary, len(policySlice))
	for i, policy := range policySlice {
		result[i] = &models.ItemSummary{
			ID:     aws.String(string(policy.ID)),
			Status: countToStatus(policy.Count),
		}
	}

	return result, nil
}

// List resources, sort by top failing
func listResources(input *dynamodb.ScanInput) ([]*models.ItemSummary, error) {
	// Fetch all resources from Dynamo
	_, resourceMap, err := scanGroupByID(input, false, true)
	if err != nil {
		return nil, err
	}

	// Convert to a slice and sort by top failing
	resourceSlice := make([]*models.ResourceSummary, 0, len(resourceMap))
	for _, resource := range resourceMap {
		resourceSlice = append(resourceSlice, resource)
	}
	sortResourcesByTopFailing(resourceSlice)

	// Convert to final ItemSummary
	result := make([]*models.ItemSummary, len(resourceSlice))
	for i, resource := range resourceSlice {
		result[i] = &models.ItemSummary{
			ID:     aws.String(string(resource.ID)),
			Status: countBySeverityToStatus(resource.Count),
		}
	}

	return result, nil
}
