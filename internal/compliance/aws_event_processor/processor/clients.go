package processor

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/panther-labs/panther/api/gateway/resources/client"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var (
	awsSession                 = session.Must(session.NewSession())
	sqsClient  sqsiface.SQSAPI = sqs.New(awsSession)
	queueURL                   = os.Getenv("SNAPSHOT_QUEUE_URL")

	transportConfig = client.DefaultTransportConfig().
			WithBasePath("/" + os.Getenv("RESOURCES_API_PATH")).
			WithHost(os.Getenv("RESOURCES_API_FQDN"))
	apiClient  = client.NewHTTPClientWithConfig(nil, transportConfig)
	httpClient = gatewayapi.GatewayClient(awsSession)
)
