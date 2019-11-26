package handlers

import (
	"errors"
	"net/http"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/compliance/models"
	"github.com/panther-labs/panther/pkg/awsbatch/dynamodbbatch"
)

// UpdateMetadata updates status entries for a given policy with a new severity / suppression set.
func UpdateMetadata(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	input, err := parseUpdateMetadata(request)
	if err != nil {
		return badRequest(err)
	}

	writes, errResponse := itemsToUpdate(input)
	if errResponse != nil {
		return errResponse
	}

	if len(writes) == 0 {
		// nothing to update
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
	}

	// It's faster to do a batch write with all of the updated entries instead of issuing
	// individual UPDATE calls for every item.
	batchInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{Env.ComplianceTable: writes},
	}

	if err := dynamodbbatch.BatchWriteItem(dynamoClient, maxWriteBackoff, batchInput); err != nil {
		zap.L().Error("dynamodbbatch.BatchWriteItem failed", zap.Error(err))
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return &events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
}

func parseUpdateMetadata(request *events.APIGatewayProxyRequest) (*models.UpdateMetadata, error) {
	var result models.UpdateMetadata
	if err := jsoniter.UnmarshalFromString(request.Body, &result); err != nil {
		return nil, err
	}

	return &result, result.Validate(nil)
}

func itemsToUpdate(input *models.UpdateMetadata) ([]*dynamodb.WriteRequest, *events.APIGatewayProxyResponse) {
	query, err := buildDescribePolicyQuery(input.PolicyID)
	if err != nil {
		return nil, &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	zap.L().Info("querying items to update",
		zap.String("policyId", string(input.PolicyID)))
	var writes []*dynamodb.WriteRequest
	err = queryPages(query, func(item *models.ComplianceStatus) error {
		ignored, patternErr := isIgnored(string(item.ResourceID), input.Suppressions)
		if patternErr != nil {
			return patternErr
		}

		// This status entry has changed - we need to rewrite it
		if bool(item.Suppressed) != ignored || item.PolicySeverity != input.Severity {
			item.PolicySeverity = input.Severity
			item.Suppressed = models.Suppressed(ignored)

			marshalled, err := dynamodbattribute.MarshalMap(item)
			if err != nil {
				return err
			}

			writes = append(writes, &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{Item: marshalled},
			})
		}

		return nil
	})

	if err != nil {
		if err == path.ErrBadPattern {
			return nil, badRequest(errors.New("invalid suppression pattern: " + err.Error()))
		}
		return nil, &events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	}

	return writes, nil
}