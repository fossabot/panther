package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/organization/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// Get retrieves account details from the table.
func (table *OrganizationsTable) Get() (*models.Organization, error) {
	zap.L().Info("retrieving organization from dynamo")
	response, err := table.client.GetItem(&dynamodb.GetItemInput{
		Key:       DynamoItem{"id": {S: aws.String("1")}},
		TableName: table.Name,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "dynamodb.GetItem", Err: err}
	}

	var org models.Organization
	if err = dynamodbattribute.UnmarshalMap(response.Item, &org); err != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to unmarshal dynamo item to an Organization: " + err.Error()}
	}
	if org.AwsConfig == nil {
		return nil, &genericapi.DoesNotExistError{}
	}

	return &org, nil
}
