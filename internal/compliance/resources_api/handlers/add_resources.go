package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/gateway/resources/models"
	"github.com/panther-labs/panther/pkg/awsbatch/dynamodbbatch"
	"github.com/panther-labs/panther/pkg/awsbatch/sqsbatch"
)

const maxBackoff = 30 * time.Second

// AddResources batch writes a group of resources to the Dynamo table.
func AddResources(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseAddResources(request)
	if err != nil {
		return badRequest(err)
	}

	now := models.LastModified(time.Now())
	writeRequests := make([]*dynamodb.WriteRequest, len(input.Resources))
	sqsEntries := make([]*sqs.SendMessageBatchRequestEntry, len(input.Resources))
	for i, r := range input.Resources {
		item := resourceItem{
			Attributes:      r.Attributes,
			Deleted:         false,
			ID:              r.ID,
			IntegrationID:   r.IntegrationID,
			IntegrationType: r.IntegrationType,
			LastModified:    now,
			Type:            r.Type,
			LowerID:         strings.ToLower(string(r.ID)),
		}

		marshalled, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			zap.L().Error("dynamodbattribute.MarshalMap failed", zap.Error(err))
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}
		writeRequests[i] = &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: marshalled}}

		body, err := jsoniter.MarshalToString(item.Resource(""))
		if err != nil {
			zap.L().Error("jsoniter.MarshalToString(resource) failed", zap.Error(err))
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}
		sqsEntries[i] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(i)),
			MessageBody: aws.String(body),
		}
	}

	dynamoInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{env.ResourcesTable: writeRequests},
	}
	if err := dynamodbbatch.BatchWriteItem(dynamoClient, maxBackoff, dynamoInput); err != nil {
		zap.L().Error("dynamodbbatch.BatchWriteItem failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	sqsInput := &sqs.SendMessageBatchInput{
		Entries:  sqsEntries,
		QueueUrl: &env.ResourcesQueueURL,
	}
	if err := sqsbatch.SendMessageBatch(sqsClient, maxBackoff, sqsInput); err != nil {
		zap.L().Error("sqsbatch.SendMessageBatch failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return &events.APIGatewayProxyResponse{StatusCode: http.StatusCreated}
}

// Parse the request body into the input model.
func parseAddResources(request *events.APIGatewayProxyRequest) (*models.AddResources, error) {
	var result models.AddResources
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	// Swagger doesn't validate plain objects
	for i, resource := range result.Resources {
		if len(resource.Attributes.(map[string]interface{})) == 0 {
			return nil, fmt.Errorf("resources[%d].attributes cannot be empty", i)
		}
	}

	return &result, result.Validate(nil)
}
