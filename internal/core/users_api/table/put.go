package users

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/pkg/genericapi"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// Put writes user to organization mapping to the dynamo table.
func (table *Table) Put(user *models.UserItem) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return &genericapi.InternalError{
			Message: "failed to marshal User to a dynamo item: " + err.Error()}
	}

	zap.L().Info("saving user to dynamo", zap.String("user email", *user.ID))
	input := &dynamodb.PutItemInput{Item: item, TableName: table.Name}
	if _, err = table.client.PutItem(input); err != nil {
		return &genericapi.AWSError{Method: "dynamodb.PutItem", Err: err}
	}

	return nil
}
