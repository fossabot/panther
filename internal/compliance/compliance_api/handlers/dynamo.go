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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/compliance/models"
)

type policyMap map[models.PolicyID]*models.PolicySummary
type resourceMap map[models.ResourceID]*models.ResourceSummary

var (
	awsSession                             = session.Must(session.NewSession())
	dynamoClient dynamodbiface.DynamoDBAPI = dynamodb.New(awsSession)
)

// Build the table key in the format Dynamo expects
func tableKey(resourceID models.ResourceID, policyID models.PolicyID) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"resourceId": {S: aws.String(string(resourceID))},
		"policyId":   {S: aws.String(string(policyID))},
	}
}

// Wrapper around dynamoClient.QueryPages that accepts a handler function to process each item.
func queryPages(input *dynamodb.QueryInput, handler func(*models.ComplianceStatus) error) error {
	var innerErr error
	err := dynamoClient.QueryPages(input, func(page *dynamodb.QueryOutput, lastPage bool) bool {
		if innerErr = handleItems(page.Items, handler); innerErr != nil {
			zap.L().Error("query handler failed", zap.Error(innerErr))
			return false // stop paging
		}
		return true
	})

	if innerErr != nil {
		return innerErr
	}
	if err != nil {
		zap.L().Error("dynamoClient.QueryPages failed", zap.Error(err))
		return err
	}

	return nil
}

// Wrapper around dynamoClient.ScanPages that accepts a handler function to process each item.
func scanPages(input *dynamodb.ScanInput, handler func(*models.ComplianceStatus) error) error {
	var innerErr error
	err := dynamoClient.ScanPages(input, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		if innerErr = handleItems(page.Items, handler); innerErr != nil {
			zap.L().Error("scan handler failed", zap.Error(innerErr))
			return false // stop paging
		}
		return true
	})

	if innerErr != nil {
		return innerErr
	}
	if err != nil {
		zap.L().Error("dynamoClient.ScanPages failed", zap.Error(err))
		return err
	}

	return nil
}

// Page handler shared by queryPages and ScanPages
func handleItems(items []map[string]*dynamodb.AttributeValue, handler func(*models.ComplianceStatus) error) error {
	var statusPage []*models.ComplianceStatus
	if err := dynamodbattribute.UnmarshalListOfMaps(items, &statusPage); err != nil {
		return err
	}

	for _, entry := range statusPage {
		if err := handler(entry); err != nil {
			return err
		}
	}

	return nil
}

// Scan Dynamo table to group everything by policyID and/or resourceID
func scanGroupByID(
	input *dynamodb.ScanInput,
	includePolicies bool,
	includeResources bool,
) (policyMap, resourceMap, error) {

	var policies policyMap
	var resources resourceMap
	if includePolicies {
		policies = make(policyMap, 200)
	}
	if includeResources {
		resources = make(resourceMap, 1000)
	}

	// Summarize every policy and resource in the organization.
	err := scanPages(input, func(item *models.ComplianceStatus) error {
		// Update policies
		if includePolicies {
			policy, ok := policies[item.PolicyID]
			if !ok {
				policy = &models.PolicySummary{
					Count:    NewStatusCount(),
					ID:       item.PolicyID,
					Severity: item.PolicySeverity,
				}
				policies[item.PolicyID] = policy
			}
			updateStatusCount(policy.Count, item.Status)
		}

		// Update resources
		if includeResources {
			resource, ok := resources[item.ResourceID]
			if !ok {
				resource = &models.ResourceSummary{
					Count: NewStatusCountBySeverity(),
					ID:    item.ResourceID,
					Type:  item.ResourceType,
				}
				resources[item.ResourceID] = resource
			}
			updateStatusCountBySeverity(resource.Count, item.PolicySeverity, item.Status)
		}

		return nil
	})

	return policies, resources, err
}
