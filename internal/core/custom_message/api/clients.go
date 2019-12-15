package custommessage

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/panther-labs/panther/internal/core/users_api/gateway"
)

var (
	appDomainURL = os.Getenv("APP_DOMAIN_URL")
	awsSession   = session.Must(session.NewSession())

	userGateway gateway.API = gateway.New(awsSession)
)
