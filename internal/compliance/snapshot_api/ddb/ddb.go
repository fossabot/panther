package ddb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	hashKey = "integrationId"
)

// DDB is a struct containing the DynamoDB client, and the table name to retrieve data.
type DDB struct {
	Client    dynamodbiface.DynamoDBAPI
	TableName string
}

// New instantiates a new client.
func New(tableName string) *DDB {
	return &DDB{
		Client:    dynamodb.New(session.Must(session.NewSession())),
		TableName: tableName,
	}
}
