package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// PutOutput saves the output details to the table.
func (table *OutputsTable) PutOutput(output *models.AlertOutputItem) error {
	item, err := dynamodbattribute.MarshalMap(output)
	if err != nil {
		return &genericapi.InternalError{Message: "failed to marshal AlertOutput to a dynamo item: " + err.Error()}
	}

	input := &dynamodb.PutItemInput{
		Item:                item,
		TableName:           table.Name,
		ConditionExpression: aws.String("attribute_not_exists(outputId)"),
	}

	if _, err = table.client.PutItem(input); err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return &genericapi.DoesNotExistError{Message: "outputId" + *output.OutputID}
		}
		return &genericapi.AWSError{Method: "dynamodb.PutItem", Err: err}
	}

	return nil
}
