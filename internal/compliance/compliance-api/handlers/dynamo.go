package handlers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/compliance/models"
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

// Wrapper around dynamoClient.ScanPages that accepts a handler function to process each item.
func queryPages(input *dynamodb.QueryInput, handler func(*models.ComplianceStatus) error) error {
	var handlerErr, unmarshalErr error

	err := dynamoClient.QueryPages(input, func(page *dynamodb.QueryOutput, lastPage bool) bool {
		var statusPage []*models.ComplianceStatus
		if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &statusPage); unmarshalErr != nil {
			return false // stop paginating
		}

		for _, entry := range statusPage {
			if handlerErr = handler(entry); handlerErr != nil {
				return false // stop paginating
			}
		}

		return true // keep paging
	})

	if handlerErr != nil {
		zap.L().Error("query item handler failed", zap.Error(handlerErr))
		return handlerErr
	}

	if unmarshalErr != nil {
		zap.L().Error("dynamodbattribute.UnmarshalListOfMaps failed", zap.Error(unmarshalErr))
		return unmarshalErr
	}

	if err != nil {
		zap.L().Error("dynamoClient.QueryPages failed", zap.Error(err))
		return err
	}

	return nil
}

// Wrapper around dynamoClient.ScanPages that accepts a handler function to process each item.
func scanPages(input *dynamodb.ScanInput, handler func(*models.ComplianceStatus) error) error {
	var handlerErr, unmarshalErr error

	err := dynamoClient.ScanPages(input, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		var statusPage []*models.ComplianceStatus
		if unmarshalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &statusPage); unmarshalErr != nil {
			return false // stop paginating
		}

		for _, entry := range statusPage {
			if handlerErr = handler(entry); handlerErr != nil {
				return false // stop paginating
			}
		}

		return true // keep paging
	})

	if handlerErr != nil {
		zap.L().Error("query item handler failed", zap.Error(handlerErr))
		return handlerErr
	}

	if unmarshalErr != nil {
		zap.L().Error("dynamodbattribute.UnmarshalListOfMaps failed", zap.Error(unmarshalErr))
		return unmarshalErr
	}

	if err != nil {
		zap.L().Error("dynamoClient.QueryPages failed", zap.Error(err))
		return err
	}

	return nil
}

// Scan Dynamo table to group everything by policyID and/or resourceID
func scanGroupByID(input *dynamodb.ScanInput, includePolicies bool, includeResources bool) (
	policies policyMap, resources resourceMap, err error) {

	if includePolicies {
		policies = make(policyMap, 200)
	}
	if includeResources {
		resources = make(resourceMap, 1000)
	}

	// Summarize every policy and resource in the organization.
	err = scanPages(input, func(item *models.ComplianceStatus) error {
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

	return
}
