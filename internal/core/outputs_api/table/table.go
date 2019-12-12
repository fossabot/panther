// Package table manages all of the Dynamo calls (query, scan, get, write, etc).
package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/panther-labs/panther/api/lambda/outputs/models"
)

// OutputsAPI defines the interface for the outputs table which can be used for mocking.
type OutputsAPI interface {
	GetOutputByName(*string) (*models.AlertOutputItem, error)
	DeleteOutput(*string) error
	GetOutputs() ([]*models.AlertOutputItem, error)
	GetOutput(*string) (*models.AlertOutputItem, error)
	PutOutput(*models.AlertOutputItem) error
	UpdateOutput(*models.AlertOutputItem) (*models.AlertOutputItem, error)
}

// OutputsTable encapsulates a connection to the Dynamo rules table.
type OutputsTable struct {
	Name             *string
	DisplayNameIndex *string
	client           dynamodbiface.DynamoDBAPI
}

// NewOutputs creates an AWS client to interface with the outputs table.
func NewOutputs(name string, displayNameIndex string, sess *session.Session) *OutputsTable {
	return &OutputsTable{
		Name:             aws.String(name),
		DisplayNameIndex: aws.String(displayNameIndex),
		client:           dynamodb.New(sess),
	}
}

// DefaultsAPI defines the interface for the table storing the default output information
type DefaultsAPI interface {
	PutDefaults(item *models.DefaultOutputsItem) error
	GetDefaults() ([]*models.DefaultOutputsItem, error)
	GetDefault(severity *string) (*models.DefaultOutputsItem, error)
}

// DefaultsTable allows interacting with DDB table storing default outputs information
type DefaultsTable struct {
	Name   *string
	client dynamodbiface.DynamoDBAPI
}

// NewDefaults creates an AWS client to interface with the defaults table.
func NewDefaults(name string, sess *session.Session) *DefaultsTable {
	return &DefaultsTable{
		Name:   aws.String(name),
		client: dynamodb.New(sess),
	}
}

// DynamoItem is a type alias for the item format expected by the Dynamo SDK.
type DynamoItem = map[string]*dynamodb.AttributeValue
