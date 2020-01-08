package handlers

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/panther-labs/panther/api/gateway/compliance/models"
)

// Common GET parameters for paging operations (DescribePolicy and DescribeResource)
type pageParams struct {
	Page       int
	PageSize   int
	Status     models.Status
	Suppressed *bool
}

// Parse page parameters for DescribePolicy and DescribeResource.
func parsePageParams(request *events.APIGatewayProxyRequest) (*pageParams, error) {
	result := pageParams{
		Page:       defaultPage,
		PageSize:   defaultPageSize,
		Status:     models.Status(request.QueryStringParameters["status"]),
		Suppressed: nil,
	}

	var err error

	if result.Status != "" {
		if err = result.Status.Validate(nil); err != nil {
			return nil, errors.New("invalid status: " + err.Error())
		}
	}

	if rawPage := request.QueryStringParameters["page"]; rawPage != "" {
		result.Page, err = strconv.Atoi(rawPage)
		if err != nil {
			return nil, errors.New("invalid page: " + err.Error())
		}
	}

	if rawPageSize := request.QueryStringParameters["pageSize"]; rawPageSize != "" {
		result.PageSize, err = strconv.Atoi(rawPageSize)
		if err != nil {
			return nil, errors.New("invalid pageSize: " + err.Error())
		}
	}

	if rawSuppressed := request.QueryStringParameters["suppressed"]; rawSuppressed != "" {
		suppressBool, err := strconv.ParseBool(rawSuppressed)
		if err != nil {
			return nil, errors.New("invalid suppressed: " + err.Error())
		}
		result.Suppressed = aws.Bool(suppressBool)
	}

	return &result, nil
}

// Common query logic for both DescribePolicy and DescribeResource.
func policyResourceDetail(
	input *dynamodb.QueryInput, params *pageParams, severity models.PolicySeverity) (*models.PolicyResourceDetail, error) {

	// TODO - global totals could be cached so not every page query has to scan everything
	result := models.PolicyResourceDetail{
		Items: make([]*models.ComplianceStatus, 0, params.PageSize),
		Paging: &models.Paging{
			TotalItems: aws.Int64(0),
		},
		Status: models.StatusPASS,
		Totals: &models.ActiveSuppressCount{
			Active:     NewStatusCount(),
			Suppressed: NewStatusCount(),
		},
	}

	err := queryPages(input, func(item *models.ComplianceStatus) error {
		addItemToResult(item, &result, params, severity)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Compute the total number of pages needed to show all the matching results
	result.Paging.TotalPages = aws.Int64(*result.Paging.TotalItems / int64(params.PageSize))
	remainder := *result.Paging.TotalItems % int64(params.PageSize)
	if remainder > 0 {
		*result.Paging.TotalPages++
	}

	if *result.Paging.TotalItems == 0 {
		result.Paging.ThisPage = aws.Int64(0)
	} else {
		result.Paging.ThisPage = aws.Int64(int64(params.Page))
	}

	return &result, nil
}

// Update the paging result with a single compliance status entry.
func addItemToResult(
	item *models.ComplianceStatus,
	result *models.PolicyResourceDetail,
	params *pageParams,
	severity models.PolicySeverity,
) {

	// Update overall status and global totals (pre-filter)
	// ERROR trumps FAIL trumps PASS
	switch item.Status {
	case models.StatusERROR:
		if item.Suppressed {
			*result.Totals.Suppressed.Error++
		} else {
			result.Status = models.StatusERROR
			*result.Totals.Active.Error++
		}

	case models.StatusFAIL:
		if item.Suppressed {
			*result.Totals.Suppressed.Fail++
		} else {
			if result.Status != models.StatusERROR {
				result.Status = models.StatusFAIL
			}
			*result.Totals.Active.Fail++
		}

	default:
		if item.Suppressed {
			*result.Totals.Suppressed.Pass++
		} else {
			*result.Totals.Active.Pass++
		}
	}

	// Drop this table entry if it doesn't match the filters
	if !itemMatchesFilter(item, params, severity) {
		return
	}

	*result.Paging.TotalItems++
	firstItem := int64((params.Page-1)*params.PageSize) + 1 // first matching item # in the requested page
	if *result.Paging.TotalItems >= firstItem && len(result.Items) < params.PageSize {
		// This matching item is in the requested page number
		result.Items = append(result.Items, item)
	}
}

func itemMatchesFilter(item *models.ComplianceStatus, params *pageParams, severity models.PolicySeverity) bool {
	if params.Suppressed != nil && *params.Suppressed != bool(item.Suppressed) {
		return false
	}
	if params.Status != "" && params.Status != item.Status {
		return false
	}
	if severity != "" && severity != item.PolicySeverity {
		return false
	}

	return true
}
