package ddb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	models "github.com/panther-labs/panther/api/snapshot"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// ScanEnabledIntegrations returns all enabled integrations based on type.
// It performs a DDB scan of the entire table with a filter expression.
func (ddb *DDB) ScanEnabledIntegrations(input *models.ListIntegrationsInput) ([]*models.SourceIntegration, error) {
	filt := expression.And(
		expression.Name("scanEnabled").Equal(expression.Value(true)),
		expression.Name("integrationType").Equal(expression.Value(input.IntegrationType)),
	)
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, &genericapi.InternalError{Message: "failed to build dynamodb expression"}
	}

	output, err := ddb.Client.Scan(&dynamodb.ScanInput{
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName:                 aws.String(ddb.TableName),
	})
	if err != nil {
		return nil, &genericapi.AWSError{Err: err, Method: "Dynamodb.Scan"}
	}

	var enabledIntegrations []*models.SourceIntegration
	if err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &enabledIntegrations); err != nil {
		return nil, err
	}

	if enabledIntegrations == nil {
		enabledIntegrations = make([]*models.SourceIntegration, 0)
	}
	return enabledIntegrations, nil
}
