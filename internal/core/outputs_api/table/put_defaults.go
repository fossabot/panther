package table

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// PutDefaults saves the default outputs to the table.
func (table *DefaultsTable) PutDefaults(defaultOutputs *models.DefaultOutputsItem) error {
	item, err := dynamodbattribute.MarshalMap(defaultOutputs)
	if err != nil {
		return &genericapi.InternalError{Message: "failed to marshal AlertOutput to a dynamo item: " + err.Error()}
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: table.Name,
	}

	if _, err = table.client.PutItem(input); err != nil {
		return &genericapi.AWSError{Method: "dynamodb.PutItem", Err: err}
	}

	return nil
}
