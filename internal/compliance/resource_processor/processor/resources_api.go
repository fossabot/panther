package processor

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
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/resources/client/operations"
)

// How many resources (with attributes) we can request in a single page.
// The goal is to keep this as high as possible while still keeping the result under 6MB.
const resourcePageSize = 2000

// Get a page of resources from the resources-api
//
// Returns {resourceID: resource}, totalPages, error
func getResources(resourceTypes []string, pageno int64) (resourceMap, int64, error) {
	result := make(resourceMap)

	zap.L().Info("listing resources from resources-api",
		zap.Int64("pageNo", pageno),
		zap.Int("pageSize", resourcePageSize),
		zap.Strings("resourceTypes", resourceTypes),
	)

	page, err := resourceClient.Operations.ListResources(&operations.ListResourcesParams{
		Deleted:    aws.Bool(false),
		Fields:     []string{"attributes", "id", "integrationId", "integrationType", "type"},
		Page:       &pageno,
		PageSize:   aws.Int64(resourcePageSize),
		Types:      resourceTypes,
		HTTPClient: httpClient,
	})
	if err != nil {
		zap.L().Error("failed to list resources", zap.Error(err))
		return nil, 0, err
	}

	for _, resource := range page.Payload.Resources {
		result[string(resource.ID)] = resource
	}
	return result, *page.Payload.Paging.TotalPages, nil
}
