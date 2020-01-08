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
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/api/gateway/analysis/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// CreatePolicy adds a new policy to the Dynamo table.
func CreatePolicy(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseUpdatePolicy(request)
	if err != nil {
		return badRequest(err)
	}

	item := &tableItem{
		AutoRemediationID:         input.AutoRemediationID,
		AutoRemediationParameters: input.AutoRemediationParameters,
		Body:                      input.Body,
		Description:               input.Description,
		DisplayName:               input.DisplayName,
		Enabled:                   input.Enabled,
		ID:                        input.ID,
		Reference:                 input.Reference,
		ResourceTypes:             input.ResourceTypes,
		Runbook:                   input.Runbook,
		Severity:                  input.Severity,
		Suppressions:              input.Suppressions,
		Tags:                      input.Tags,
		Tests:                     input.Tests,
		Type:                      typePolicy,
	}

	if _, err := writeItem(item, input.UserID, aws.Bool(false)); err != nil {
		if err == errExists {
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusConflict}
		}
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	// New policies are "passing" since they haven't evaluated anything yet.
	return gatewayapi.MarshalResponse(item.Policy(models.ComplianceStatusPASS), http.StatusCreated)
}

// body parsing shared by CreatePolicy and ModifyPolicy
func parseUpdatePolicy(request *events.APIGatewayProxyRequest) (*models.UpdatePolicy, error) {
	var result models.UpdatePolicy
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	if err := result.Validate(nil); err != nil {
		return nil, err
	}

	return &result, nil
}
