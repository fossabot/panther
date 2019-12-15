package users

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/panther-labs/panther/api/lambda/users/models"
)

// API defines the interface for the table which can be used for mocking.
type API interface {
	Delete(*string) error
	Get(*string) (*models.UserItem, error)
	Put(user *models.UserItem) error
}

// Table encapsulates a connection to the Dynamo table.
type Table struct {
	Name   *string
	client dynamodbiface.DynamoDBAPI
}

// The Table must satisfy the API interface.
var _ API = (*Table)(nil)

// New creates a new Dynamo client which talks to the given table name.
func New(tableName string, sess *session.Session) *Table {
	return &Table{Name: aws.String(tableName), client: dynamodb.New(sess)}
}

// DynamoItem is a type alias for the item format expected by the Dynamo SDK.
type DynamoItem = map[string]*dynamodb.AttributeValue
