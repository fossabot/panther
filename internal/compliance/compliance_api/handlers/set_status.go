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
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/compliance/models"
	"github.com/panther-labs/panther/pkg/awsbatch/dynamodbbatch"
)

const (
	maxWriteBackoff = time.Minute

	// Automatically expire status entries if they haven't been updated in 2 days.
	//
	// This handles edge cases where deleted policies/resources aren't fully cleared due to
	// eventual consistency, queue delays, etc.
	statusLifetime = 50 * time.Hour
)

// SetStatus batch writes a set of compliance status to the Dynamo table.
func SetStatus(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseSetStatus(request)
	if err != nil {
		return badRequest(err)
	}

	now := time.Now()
	expiresAt := now.Add(statusLifetime).Unix()
	writeRequests := make([]*dynamodb.WriteRequest, len(input.Entries))
	for i, entry := range input.Entries {
		status := &models.ComplianceStatus{
			ErrorMessage:   entry.ErrorMessage,
			ExpiresAt:      models.ExpiresAt(expiresAt),
			IntegrationID:  entry.IntegrationID,
			LastUpdated:    models.LastUpdated(now),
			PolicyID:       entry.PolicyID,
			PolicySeverity: entry.PolicySeverity,
			ResourceID:     entry.ResourceID,
			ResourceType:   entry.ResourceType,
			Status:         entry.Status,
			Suppressed:     entry.Suppressed,
		}

		marshalled, err := dynamodbattribute.MarshalMap(status)
		if err != nil {
			zap.L().Error("dynamodbattribute.MarshalMap failed", zap.Error(err))
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}

		writeRequests[i] = &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: marshalled}}
	}

	batchInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{Env.ComplianceTable: writeRequests},
	}

	if err := dynamodbbatch.BatchWriteItem(dynamoClient, maxWriteBackoff, batchInput); err != nil {
		zap.L().Error("dynamodbbatch.BatchWriteItem failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return &events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}
}

func parseSetStatus(request *events.APIGatewayProxyRequest) (*models.SetStatusBatch, error) {
	var result models.SetStatusBatch
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	return &result, result.Validate(nil)
}
