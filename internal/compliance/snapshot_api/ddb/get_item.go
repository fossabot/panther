package ddb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// GetIntegration returns an integration by its ID
func (ddb *DDB) GetIntegration(integrationID *string) (*models.SourceIntegrationMetadata, error) {
	output, err := ddb.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ddb.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			hashKey: {S: integrationID},
		},
	})
	if err != nil {
		return nil, &genericapi.AWSError{Err: err, Method: "Dynamodb.GetItem"}
	}

	var integration models.SourceIntegrationMetadata
	if output.Item == nil {
		return nil, err
	}
	if err := dynamodbattribute.UnmarshalMap(output.Item, &integration); err != nil {
		return nil, err
	}

	return &integration, nil
}
