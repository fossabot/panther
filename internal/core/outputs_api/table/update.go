package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// UpdateOutput updates an existing item in the table
func (table *OutputsTable) UpdateOutput(alertOutput *models.AlertOutputItem) (*models.AlertOutputItem, error) {
	updateExpression := expression.
		Set(expression.Name("displayName"), expression.Value(alertOutput.DisplayName)).
		Set(expression.Name("lastModifiedBy"), expression.Value(alertOutput.LastModifiedBy)).
		Set(expression.Name("lastModifiedTime"), expression.Value(alertOutput.LastModifiedTime)).
		Set(expression.Name("outputType"), expression.Value(alertOutput.OutputType)).
		Set(expression.Name("encryptedConfig"), expression.Value(alertOutput.EncryptedConfig))

	if alertOutput.VerificationStatus != nil {
		updateExpression.Set(expression.Name("verificationStatus"), expression.Value(alertOutput.VerificationStatus))
	}

	conditionExpression := expression.Name("outputId").Equal(expression.Value(alertOutput.OutputID))
	combinedExpression, err := expression.NewBuilder().
		WithCondition(conditionExpression).
		WithUpdate(updateExpression).
		Build()

	if err != nil {
		return nil, &genericapi.InternalError{Message: "failed to build expression " + err.Error()}
	}

	updateResult, err := table.client.UpdateItem(
		&dynamodb.UpdateItemInput{
			TableName: table.Name,
			Key: DynamoItem{
				"outputId": {S: alertOutput.OutputID},
			},
			UpdateExpression:          combinedExpression.Update(),
			ConditionExpression:       combinedExpression.Condition(),
			ExpressionAttributeNames:  combinedExpression.Names(),
			ExpressionAttributeValues: combinedExpression.Values(),
			ReturnValues:              aws.String(dynamodb.ReturnValueAllNew),
		})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return nil, &genericapi.DoesNotExistError{Message: "outputId=" + *alertOutput.OutputID}
		}
		return nil, &genericapi.AWSError{Method: "dynamodb.UpdateItem", Err: err}
	}

	var output models.AlertOutputItem
	if err = dynamodbattribute.UnmarshalMap(updateResult.Attributes, &output); err != nil {
		return nil, &genericapi.InternalError{
			Message: "failed to unmarshal dynamo item to an AlertOutputItem: " + err.Error()}
	}
	return &output, nil
}
