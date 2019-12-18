// Package table manages all of the Dynamo calls (query, scan, get, write, etc).
package table

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/panther-labs/panther/api/lambda/alerts/models"
)

// API defines the interface for the alerts table which can be used for mocking.
type API interface {
	GetAlert(*string) (*models.AlertItem, error)
	GetEvent([]byte) (*string, error)
	ListAlerts(*string, *int) ([]*models.AlertItem, *string, error)
	ListAlertsByRule(*string, *string, *int) ([]*models.AlertItem, *string, error)
}

// AlertsTable encapsulates a connection to the Dynamo alerts table.
type AlertsTable struct {
	AlertsTableName             string
	RuleIDCreationTimeIndexName string
	EventsTableName             string
	Client                      dynamodbiface.DynamoDBAPI
}

// The AlertsTable must satisfy the API interface.
var _ API = (*AlertsTable)(nil)

// DynamoItem is a type alias for the item format expected by the Dynamo SDK.
type DynamoItem = map[string]*dynamodb.AttributeValue
