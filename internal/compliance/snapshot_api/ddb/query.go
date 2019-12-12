package ddb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	models "github.com/panther-labs/panther/api/snapshot"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// QueryIntegrations returns all snapshot integrations
func (ddb *DDB) QueryIntegrations() ([]*models.SourceIntegration, error) {
	builder := expression.NewBuilder()
	expr, err := builder.Build()
	if err != nil {
		return nil, &genericapi.InternalError{Message: "failed to build dynamodb expression"}
	}

	output, err := ddb.Client.Query(&dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(ddb.TableName),
	})
	if err != nil {
		return nil, &genericapi.AWSError{Err: err, Method: "Dynamodb.Scan"}
	}

	var integrations []*models.SourceIntegration
	if err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &integrations); err != nil {
		return nil, err
	}

	if integrations == nil {
		integrations = make([]*models.SourceIntegration, 0)
	}

	return integrations, nil
}
