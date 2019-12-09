package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/kelseyhightower/envconfig"

	complianceapi "github.com/panther-labs/panther/api/compliance/client"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var (
	env envConfig

	awsSession   *session.Session
	dynamoClient dynamodbiface.DynamoDBAPI
	s3Client     s3iface.S3API
	sqsClient    sqsiface.SQSAPI

	httpClient       *http.Client
	complianceClient *complianceapi.PantherCompliance
)

type envConfig struct {
	Bucket            string `required:"true" split_words:"true"`
	ComplianceAPIHost string `required:"true" split_words:"true"`
	ComplianceAPIPath string `required:"true" split_words:"true"`
	Engine            string `required:"true" split_words:"true"`
	ResourceQueueURL  string `required:"true" split_words:"true"`
	Table             string `required:"true" split_words:"true"`
}

// Setup parses the environment and constructs AWS and http clients on a cold Lambda start.
// All required environment variables must be present or this function will panic.
func Setup() {
	envconfig.MustProcess("", &env)

	awsSession = session.Must(session.NewSession())
	dynamoClient = dynamodb.New(awsSession)
	s3Client = s3.New(awsSession)
	sqsClient = sqs.New(awsSession)

	httpClient = gatewayapi.GatewayClient(awsSession)
	complianceClient = complianceapi.NewHTTPClientWithConfig(
		nil, complianceapi.DefaultTransportConfig().
			WithBasePath("/"+env.ComplianceAPIPath).
			WithHost(env.ComplianceAPIHost))
}
