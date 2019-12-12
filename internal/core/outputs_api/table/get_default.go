package table

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// GetDefault gets the default outputs for a given severity
func (table *DefaultsTable) GetDefault(severity *string) (*models.DefaultOutputsItem, error) {
	result, err := table.client.GetItem(
		&dynamodb.GetItemInput{
			TableName: table.Name,
			Key: DynamoItem{
				"severity": {S: severity},
			},
		})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, err
	}
	var defaultOutput models.DefaultOutputsItem
	if err = dynamodbattribute.UnmarshalMap(result.Item, &defaultOutput); err != nil {
		return nil, err
	}
	return &defaultOutput, nil
}
