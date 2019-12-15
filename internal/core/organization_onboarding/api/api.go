// Package api defines CRUD actions for the Cognito Api.
package api

import (
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/panther-labs/panther/internal/core/organization_onboarding/gateway"
)

// The API has receiver methods for each of the handlers.
type API struct{}

var (
	awsSession                      = session.Must(session.NewSession())
	stepFunctionGateway gateway.API = gateway.New(awsSession)
)
