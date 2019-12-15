package users

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// Delete removes a row from the table.
func (table *Table) Delete(id *string) error {
	condition := expression.AttributeExists(expression.Name("id"))
	expr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		return &genericapi.InternalError{Message: "dynamo expression build failed: " + err.Error()}
	}

	zap.L().Info("deleting user from dynamo", zap.String("id", *id))
	_, err = table.client.DeleteItem(&dynamodb.DeleteItemInput{
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       DynamoItem{"id": {S: id}},
		TableName:                 table.Name,
	})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return &genericapi.DoesNotExistError{Message: "id=" + *id}
		}
		return &genericapi.AWSError{Method: "dynamodb.DeleteItem", Err: err}
	}

	return nil
}
