package apihandlers

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
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/pkg/gatewayapi"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// GetRemediations returns the list of remediations available for an organization
func GetRemediations(_ *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	zap.L().Info("getting list of remediations")
	// TODO - differentiate between different error types
	remediations, err := invoker.GetRemediations()
	if err != nil {
		if _, ok := err.(*genericapi.DoesNotExistError); ok {
			return gatewayapi.MarshalResponse(RemediationLambdaNotFound, http.StatusNotFound)
		}
		zap.L().Warn("failed to fetch available remediations", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	body, err := jsoniter.MarshalToString(remediations)
	if err != nil {
		zap.L().Error("failed to marshal remediations", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: body}
}
