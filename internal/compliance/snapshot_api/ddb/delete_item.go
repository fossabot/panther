package ddb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	models "github.com/panther-labs/panther/api/snapshot"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// DeleteIntegrationItem deletes an integration from the database based on the integration ID
func (ddb *DDB) DeleteIntegrationItem(input *models.DeleteIntegrationInput) error {
	condition := expression.AttributeExists(expression.Name("integrationId"))

	builder := expression.NewBuilder().WithCondition(condition)
	expr, err := builder.Build()
	if err != nil {
		return &genericapi.InternalError{Message: "failed to build DeleteIntegration ddb expression"}
	}

	_, err = ddb.Client.DeleteItem(&dynamodb.DeleteItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		Key: map[string]*dynamodb.AttributeValue{
			hashKey: {S: input.IntegrationID},
		},
		TableName: aws.String(ddb.TableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return &genericapi.DoesNotExistError{Message: aerr.Error()}
			default:
				return &genericapi.AWSError{Err: err, Method: "Dynamodb.DeleteItem"}
			}
		}
	}

	return nil
}
