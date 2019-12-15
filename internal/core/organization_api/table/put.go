package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// Put writes organization details to the dynamo table.
func (table *OrganizationsTable) Put(org *models.Organization) error {
	item, err := dynamodbattribute.MarshalMap(org)
	if err != nil {
		return &genericapi.InternalError{
			Message: "failed to marshal Organization to a dynamo item: " + err.Error()}
	}

	item["id"] = &dynamodb.AttributeValue{S: aws.String("1")}
	input := &dynamodb.PutItemInput{Item: item, TableName: table.Name}
	if _, err = table.client.PutItem(input); err != nil {
		return &genericapi.AWSError{Method: "dynamodb.PutItem", Err: err}
	}

	return nil
}
