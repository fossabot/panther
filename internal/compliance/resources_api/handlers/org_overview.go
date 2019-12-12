package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/resources/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// OrgOverview retrieves a summary of the resources that exist in an organization
func OrgOverview(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	scanInput, err := buildOrgOverviewScan()
	if err != nil {
		zap.L().Error("failed to build query", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	orgOverviewMap := make(map[models.ResourceType]int64)
	err = scanPages(scanInput, func(item *resourceItem) error {
		orgOverviewMap[item.Type]++
		return nil
	})

	if err != nil {
		zap.L().Warn("failed to query dynamo", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	orgOverview := &models.OrgOverview{
		Resources: make([]*models.ResourceTypeSummary, 0, len(orgOverviewMap)),
	}
	for resourceType, count := range orgOverviewMap {
		summary := &models.ResourceTypeSummary{Count: aws.Int64(count), Type: resourceType}
		orgOverview.Resources = append(orgOverview.Resources, summary)
	}

	return gatewayapi.MarshalResponse(orgOverview, http.StatusOK)
}

// Building the query for getting all resources of an organization
func buildOrgOverviewScan() (*dynamodb.ScanInput, error) {
	filter := expression.Equal(expression.Name("deleted"), expression.Value(false))
	projection := expression.NamesList(expression.Name("type"))
	expr, err := expression.NewBuilder().
		WithFilter(filter).
		WithProjection(projection).
		Build()
	if err != nil {
		return nil, err
	}

	return &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 &env.ResourcesTable,
	}, nil
}
