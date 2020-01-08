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
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
)

// Suppress adds suppressions for one or more policies in the same organization.
func Suppress(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseSuppress(request)
	if err != nil {
		return badRequest(err)
	}

	updates, err := addSuppressions(input.PolicyIds, input.ResourcePatterns)
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	// Update compliance status with new suppressions
	for _, policy := range updates {
		if err := updateComplianceMetadata(policy); err != nil {
			// Log an error, but don't mark the API call as a failure
			zap.L().Error("failed to update compliance entries with new suppression", zap.Error(err))
		}
	}

	return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

func parseSuppress(request *events.APIGatewayProxyRequest) (*models.Suppress, error) {
	var result models.Suppress
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	if err := result.Validate(nil); err != nil {
		return nil, err
	}

	if len(result.ResourcePatterns) == 0 {
		return nil, errors.New("invalid resourcePatterns: at least one is required")
	}

	return &result, nil
}
