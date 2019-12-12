package ddb

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	models "github.com/panther-labs/panther/api/snapshot"
	"github.com/panther-labs/panther/pkg/awsbatch/dynamodbbatch"
	"github.com/panther-labs/panther/pkg/genericapi"
)

var maxElapsedTime = 15 * time.Second

// BatchPutSourceIntegrations adds a batch of new Snapshot Integrations to the database.
func (ddb *DDB) BatchPutSourceIntegrations(input []*models.SourceIntegrationMetadata) error {
	writeRequests := make([]*dynamodb.WriteRequest, len(input))

	// Marshal each new integration and add to the write request
	for i, item := range input {
		item, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			return &genericapi.AWSError{Err: err, Method: "Dynamodb.MarshalMap"}
		}
		writeRequests[i] = &dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: item}}
	}

	// Do the batch write
	err := dynamodbbatch.BatchWriteItem(ddb.Client, maxElapsedTime, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{ddb.TableName: writeRequests}})
	if err != nil {
		return &genericapi.AWSError{Err: err, Method: "Dynamodb.BatchWriteItem"}
	}

	return nil
}
