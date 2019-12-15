// Package api defines CRUD actions for the Cognito Api.
package api

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"
	users "github.com/panther-labs/panther/internal/core/users_api/table"
)

// The API has receiver methods for each of the handlers.
type API struct{}

var (
	organizationAPI = os.Getenv("ORGANIZATION_API")
	awsSession      = session.Must(session.NewSession())

	lambdaClient lambdaiface.LambdaAPI = lambda.New(awsSession)
	userGateway  gateway.API           = gateway.New(awsSession)
	userTable    users.API             = users.New(os.Getenv("USERS_TABLE_NAME"), awsSession)
)
