package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/compliance/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

type getOrgOverviewParams struct {
	LimitTopFailing int
}

// GetOrgOverview returns all the pass/fail information for the Panther overview dashboard.
func GetOrgOverview(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	params, err := parseGetOrgOverview(request)
	if err != nil {
		return badRequest(err)
	}

	input, err := buildGetOrgOverviewQuery()
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	policies, resources, err := scanGroupByID(input, true, true)
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(
		buildOverview(policies, resources, params.LimitTopFailing), http.StatusOK)
}

func parseGetOrgOverview(request *events.APIGatewayProxyRequest) (*getOrgOverviewParams, error) {
	result := getOrgOverviewParams{
		LimitTopFailing: defaultTopFailing,
	}

	var err error
	rawTopFailing := request.QueryStringParameters["limitTopFailing"]
	if rawTopFailing != "" {
		result.LimitTopFailing, err = strconv.Atoi(rawTopFailing)
		if err != nil {
			return nil, errors.New("invalid limitTopFailing: " + err.Error())
		}
	}

	return &result, nil
}

func buildGetOrgOverviewQuery() (*dynamodb.ScanInput, error) {
	filter := expression.Equal(expression.Name("suppressed"), expression.Value(false))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		zap.L().Error("expression.Build failed", zap.Error(err))
		return nil, err
	}

	return &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 &Env.ComplianceTable,
	}, nil
}

func buildOverview(policies policyMap, resources resourceMap, limitTopFailing int) *models.OrgSummary {
	// Count policies by severity and record failed policies
	appliedPolicies := NewStatusCountBySeverity()
	failedPolicies := make([]*models.PolicySummary, 0)
	for _, policy := range policies {
		status := countToStatus(policy.Count)
		updateStatusCountBySeverity(appliedPolicies, policy.Severity, status)
		if status != models.StatusPASS {
			failedPolicies = append(failedPolicies, policy)
		}
	}

	// Sort and truncate failed policies
	sortPoliciesByTopFailing(failedPolicies)
	if len(failedPolicies) > limitTopFailing {
		failedPolicies = failedPolicies[:limitTopFailing]
	}

	// Count resources by type and record failed resources
	resourcesByType := make(map[models.ResourceType]*models.StatusCount, 100)
	failedResources := make([]*models.ResourceSummary, 0, len(resources)/2)
	for _, resource := range resources {
		count, ok := resourcesByType[resource.Type]
		if !ok {
			count = NewStatusCount()
			resourcesByType[resource.Type] = count
		}

		status := countBySeverityToStatus(resource.Count)
		updateStatusCount(count, status)
		if status != models.StatusPASS {
			failedResources = append(failedResources, resource)
		}
	}

	// Convert resourcesByType into appropriate struct
	scannedResources := &models.ScannedResources{
		ByType: make([]*models.ResourceOfType, 0, len(resourcesByType)),
	}
	for resourceType, count := range resourcesByType {
		entry := &models.ResourceOfType{Count: count, Type: resourceType}
		scannedResources.ByType = append(scannedResources.ByType, entry)
	}

	// Sort and truncate failing resources
	sortResourcesByTopFailing(failedResources)
	if len(failedResources) > limitTopFailing {
		failedResources = failedResources[:limitTopFailing]
	}

	return &models.OrgSummary{
		AppliedPolicies:     appliedPolicies,
		ScannedResources:    scannedResources,
		TopFailingPolicies:  failedPolicies,
		TopFailingResources: failedResources,
	}
}
