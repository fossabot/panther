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

	"github.com/panther-labs/panther/api/gateway/resources/models"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// GetResource retrieves a single resource from the Dynamo table.
func GetResource(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	resourceID, err := parseGetResource(request)
	if err != nil {
		return badRequest(err)
	}

	response, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		Key:       tableKey(resourceID),
		TableName: &env.ResourcesTable,
	})
	if err != nil {
		zap.L().Error("dynamoClient.GetItem failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	if len(response.Item) == 0 {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}
	}

	var item resourceItem
	if err := dynamodbattribute.UnmarshalMap(response.Item, &item); err != nil {
		zap.L().Error("dynamodbattribute.UnmarshalMap failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	status, err := getComplianceStatus(resourceID)
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return gatewayapi.MarshalResponse(item.Resource(status.Status), http.StatusOK)
}

// API gateway doesn't do advanced validation of query parameters, but we can do it here.
func parseGetResource(request *events.APIGatewayProxyRequest) (resourceID models.ResourceID, err error) {
	escaped, err := url.QueryUnescape(request.QueryStringParameters["resourceId"])
	if err != nil {
		err = errors.New("invalid resourceId: " + err.Error())
		return
	}

	resourceID = models.ResourceID(escaped)
	if err = resourceID.Validate(nil); err != nil {
		err = errors.New("invalid resourceId: " + err.Error())
	}
	return
}
