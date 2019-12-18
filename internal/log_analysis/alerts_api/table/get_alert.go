package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
)

// GetAlert retrieve a AlertItem from DDB
func (table *AlertsTable) GetAlert(alertID *string) (*models.AlertItem, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"alertId": {S: alertID},
		},
		TableName: aws.String(table.AlertsTableName),
	}

	ddbResult, err := table.Client.GetItem(input)
	if err != nil {
		return nil, err
	}

	alertItem := &models.AlertItem{}
	if err = dynamodbattribute.UnmarshalMap(ddbResult.Item, alertItem); err != nil {
		return nil, err
	}
	return alertItem, nil
}
