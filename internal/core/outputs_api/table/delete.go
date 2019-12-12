package table

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/panther-labs/panther/pkg/genericapi"
)

// DeleteOutput removes an output from the table.
func (table *OutputsTable) DeleteOutput(outputID *string) error {
	condition := expression.Name("outputId").Equal(expression.Value(outputID))

	conditionExpression, err := expression.NewBuilder().WithCondition(condition).Build()

	if err != nil {
		return &genericapi.InternalError{Message: "failed to build expression " + err.Error()}
	}

	_, err = table.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: table.Name,
		Key: DynamoItem{
			"outputId": {S: outputID},
		},
		ConditionExpression:       conditionExpression.Condition(),
		ExpressionAttributeNames:  conditionExpression.Names(),
		ExpressionAttributeValues: conditionExpression.Values(),
	})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return &genericapi.DoesNotExistError{Message: "outputId=" + *outputID}
		}
		return &genericapi.AWSError{Method: "dynamodb.DeleteItem", Err: err}
	}

	return nil
}
