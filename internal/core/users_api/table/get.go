package users

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/users/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// Get retrieves user to org mapping from the table.
func (table *Table) Get(id *string) (*models.UserItem, error) {
	zap.L().Info("retrieving user from dynamo", zap.String("id", *id))
	response, err := table.client.GetItem(&dynamodb.GetItemInput{
		Key:       DynamoItem{"id": {S: id}},
		TableName: table.Name,
	})
	if err != nil {
		return nil, &genericapi.AWSError{Method: "dynamodb.GetItem", Err: err}
	}

	var user models.UserItem
	if err = dynamodbattribute.UnmarshalMap(response.Item, &user); err != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to unmarshal dynamo item to a User: " + err.Error()}
	}

	if aws.StringValue(user.ID) == "" {
		return nil, &genericapi.DoesNotExistError{Message: "id=" + *id}
	}

	return &user, nil
}
