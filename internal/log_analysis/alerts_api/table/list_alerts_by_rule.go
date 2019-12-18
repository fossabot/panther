package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ListAlertsByRule returns (a page of alerts, last evaluated key, any error)
func (table *AlertsTable) ListAlertsByRule(ruleID *string, exclusiveStartKey *string, pageSize *int) (
	summaries []*models.AlertItem, lastEvaluatedKey *string, err error) {

	keyCondition := expression.Key("ruleId").Equal(expression.Value(ruleID))

	queryExpression, err := expression.NewBuilder().
		WithKeyCondition(keyCondition).
		Build()

	if err != nil {
		return nil, nil, &genericapi.InternalError{Message: "failed to build expression " + err.Error()}
	}

	var queryResultsLimit *int64
	if pageSize != nil {
		queryResultsLimit = aws.Int64(int64(*pageSize))
	}

	var queryExclusiveStartKey map[string]*dynamodb.AttributeValue
	if exclusiveStartKey != nil {
		key := &listAlertsPaginationKey{}
		err = jsoniter.UnmarshalFromString(*exclusiveStartKey, key)
		if err != nil {
			return nil, nil, err
		}
		queryExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"creationTime": {S: key.CreationTime},
			"alertId":      {S: key.AlertID},
		}
	}

	var queryInput = &dynamodb.QueryInput{
		TableName:                 &table.AlertsTableName,
		ScanIndexForward:          aws.Bool(false),
		ExpressionAttributeNames:  queryExpression.Names(),
		ExpressionAttributeValues: queryExpression.Values(),
		KeyConditionExpression:    queryExpression.KeyCondition(),
		ExclusiveStartKey:         queryExclusiveStartKey,
		IndexName:                 aws.String(table.RuleIDCreationTimeIndexName),
		Limit:                     queryResultsLimit,
	}

	queryOutput, err := table.Client.Query(queryInput)
	if err != nil {
		return nil, nil, &genericapi.InternalError{
			Message: "query to DDB failed" + err.Error(),
		}
	}

	err = dynamodbattribute.UnmarshalListOfMaps(queryOutput.Items, &summaries)
	if err != nil {
		return nil, nil, &genericapi.InternalError{
			Message: "failed to marshall response" + err.Error(),
		}
	}

	// If DDB returned a LastEvaluatedKey, it means there are more alerts to be returned
	// Return populated `lastEvaluatedKey` in the response.
	if len(queryOutput.LastEvaluatedKey) > 0 {
		paginationKey := listAlertsPaginationKey{
			CreationTime: queryOutput.LastEvaluatedKey["creationTime"].S,
			AlertID:      queryOutput.LastEvaluatedKey["alertId"].S,
		}
		marshalledKey, err := jsoniter.MarshalToString(paginationKey)
		if err != nil {
			return nil, nil, &genericapi.InternalError{
				Message: "failed to marshall key" + err.Error(),
			}
		}
		lastEvaluatedKey = &marshalledKey
	}

	return summaries, lastEvaluatedKey, nil
}
