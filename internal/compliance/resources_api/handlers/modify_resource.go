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
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/api/gateway/resources/models"
)

// ModifyResource will update some of the resource properties.
func ModifyResource(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseModifyResource(request)
	if err != nil {
		return badRequest(err)
	}

	// Replace subsets of the resource attributes, not the whole thing.
	// For example, {"VersioningEnabled": true, "EncryptionConfig.KeyID": "abc"}
	update := expression.Set(expression.Name("lastModified"), expression.Value(time.Now()))
	for key, val := range input.ReplaceAttributes.(map[string]interface{}) {
		update = update.Set(expression.Name("attributes."+key), expression.Value(val))
	}

	return doUpdate(update, input.ID)
}

// Parse the request body into a ModifyResource model.
func parseModifyResource(request *events.APIGatewayProxyRequest) (*models.ModifyResource, error) {
	var result models.ModifyResource
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	// swagger doesn't validate an arbitrary object
	if len(result.ReplaceAttributes.(map[string]interface{})) == 0 {
		return &result, errors.New("at least one attribute is required")
	}

	return &result, result.Validate(nil)
}
