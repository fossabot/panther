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

	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// ModifyRule updates an existing rule.
func ModifyRule(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseUpdateRule(request)
	if err != nil {
		return badRequest(err)
	}

	item := &tableItem{
		Body:          input.Body,
		Description:   input.Description,
		DisplayName:   input.DisplayName,
		Enabled:       input.Enabled,
		ID:            input.ID,
		Reference:     input.Reference,
		ResourceTypes: input.LogTypes,
		Runbook:       input.Runbook,
		Severity:      input.Severity,
		Tags:          input.Tags,
		Tests:         input.Tests,
		Type:          typeRule,
	}

	if _, err := writeItem(item, input.UserID, aws.Bool(true)); err != nil {
		if err == errNotExists || err == errWrongType {
			// errWrongType means we tried to modify a rule which is actually a policy.
			// In this case return 404 - the rule you tried to modify does not exist.
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}
		}
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(item.Rule(), http.StatusOK)
}
