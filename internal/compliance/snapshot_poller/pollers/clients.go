package pollers

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/panther-labs/panther/api/resources/client"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var (
	transportConfig = client.DefaultTransportConfig().
			WithBasePath("/" + os.Getenv("RESOURCES_API_PATH")).
			WithHost(os.Getenv("RESOURCES_API_FQDN"))
	apiClient  = client.NewHTTPClientWithConfig(nil, transportConfig)
	awsSession = session.Must(session.NewSession())
	httpClient = gatewayapi.GatewayClient(awsSession)
)
