package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/kelseyhightower/envconfig"

	complianceapi "github.com/panther-labs/panther/api/gateway/compliance/client"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var (
	env envConfig

	awsSession   *session.Session
	dynamoClient dynamodbiface.DynamoDBAPI
	sqsClient    sqsiface.SQSAPI

	httpClient       *http.Client
	complianceClient *complianceapi.PantherCompliance
)

type envConfig struct {
	ComplianceAPIHost string `required:"true" split_words:"true"`
	ComplianceAPIPath string `required:"true" split_words:"true"`
	ResourcesQueueURL string `required:"true" split_words:"true"`
	ResourcesTable    string `required:"true" split_words:"true"`
}

// Setup parses the environment and builds the AWS and http clients.
func Setup() {
	envconfig.MustProcess("", &env)

	awsSession = session.Must(session.NewSession())
	dynamoClient = dynamodb.New(awsSession)
	sqsClient = sqs.New(awsSession)

	httpClient = gatewayapi.GatewayClient(awsSession)
	complianceClient = complianceapi.NewHTTPClientWithConfig(
		nil, complianceapi.DefaultTransportConfig().
			WithHost(env.ComplianceAPIHost).
			WithBasePath("/"+env.ComplianceAPIPath))
}
