package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetEvent retrieves an event from DDB
func (table *AlertsTable) GetEvent(eventHash []byte) (*string, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"eventHash": {B: eventHash},
		},
		TableName: aws.String(table.EventsTableName),
	}

	ddbResult, err := table.Client.GetItem(input)
	if err != nil {
		return nil, err
	}

	return ddbResult.Item["event"].S, nil
}
