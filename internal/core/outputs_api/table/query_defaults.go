package table

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// GetDefaults returns the default outputs for one organization
func (table *DefaultsTable) GetDefaults() (defaults []*models.DefaultOutputsItem, err error) {
	var scanInput = &dynamodb.ScanInput{
		TableName: table.Name,
	}

	var internalErr error
	queryErr := table.client.ScanPages(scanInput,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			var defaultsPartial []*models.DefaultOutputsItem
			if internalErr = dynamodbattribute.UnmarshalListOfMaps(page.Items, &defaultsPartial); internalErr != nil {
				internalErr = &genericapi.InternalError{
					Message: "failed to unmarshal dynamo item to an AlertOutputItem: " + internalErr.Error(),
				}
				return false
			}
			defaults = append(defaults, defaultsPartial...)
			return true
		})

	if queryErr != nil {
		return nil, &genericapi.AWSError{Err: queryErr, Method: "dynamodb.ScanPages"}
	}
	if internalErr != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to unmarshal dynamo items: " + internalErr.Error()}
	}

	return defaults, nil
}
