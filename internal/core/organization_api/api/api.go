// Package api defines CRUD actions for the Panther organization database.
package api

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/panther-labs/panther/internal/core/organization_api/table"
)

var (
	awsSession           = session.Must(session.NewSession())
	orgTable   table.API = table.New(os.Getenv("ORG_TABLE_NAME"), awsSession)
)

// API has all of the handlers as receiver methods.
type API struct{}
