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
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/compliance/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

type getParams struct {
	ResourceID models.ResourceID
	PolicyID   models.PolicyID
}

// GetStatus retrieves a single policy/resource status pair from the Dynamo table.
func GetStatus(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseGetStatus(request)
	if err != nil {
		return badRequest(err)
	}

	response, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		Key:       tableKey(input.ResourceID, input.PolicyID),
		TableName: &Env.ComplianceTable,
	})
	if err != nil {
		zap.L().Error("dynamoClient.GetItem failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	if len(response.Item) == 0 {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}
	}

	var status models.ComplianceStatus
	if err := dynamodbattribute.UnmarshalMap(response.Item, &status); err != nil {
		zap.L().Error("dynamodbattribute.UnmarshalMap failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(&status, http.StatusOK)
}

func parseGetStatus(request *events.APIGatewayProxyRequest) (*getParams, error) {
	escaped, err := url.QueryUnescape(request.QueryStringParameters["resourceId"])
	if err != nil {
		return nil, errors.New("invalid resourceId: query unescape: " + err.Error())
	}
	resourceID := models.ResourceID(escaped)
	if err = resourceID.Validate(nil); err != nil {
		return nil, errors.New("invalid resourceId: " + err.Error())
	}

	escaped, err = url.QueryUnescape(request.QueryStringParameters["policyId"])
	if err != nil {
		return nil, errors.New("invalid policyId: query unescape: " + err.Error())
	}
	policyID := models.PolicyID(escaped)
	if err = policyID.Validate(nil); err != nil {
		return nil, errors.New("invalid policyId: " + err.Error())
	}

	return &getParams{
		ResourceID: resourceID,
		PolicyID:   policyID,
	}, nil
}
