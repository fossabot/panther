// Package api defines CRUD actions for the Panther alerts database.
package api

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kelseyhightower/envconfig"

	"github.com/panther-labs/panther/api/gateway/analysis/client"
	"github.com/panther-labs/panther/internal/log_analysis/alerts_api/table"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

// API has all of the handlers as receiver methods.
type API struct{}

var (
	env            envConfig
	awsSession     *session.Session
	httpClient     *http.Client
	policiesClient *client.PantherAnalysis
	alertsDB       table.API
)

type envConfig struct {
	AnalysisAPIHost string `required:"true" split_words:"true"`
	AnalysisAPIPath string `required:"true" split_words:"true"`
	AlertsTableName string `required:"true" split_words:"true"`
	RuleIndexName   string `required:"true" split_words:"true"`
	EventsTableName string `required:"true" split_words:"true"`
}

// Setup parses the environment and builds the AWS and http clients.
func Setup() {
	envconfig.MustProcess("", &env)

	awsSession = session.Must(session.NewSession())

	httpClient = gatewayapi.GatewayClient(awsSession)
	policiesClient = client.NewHTTPClientWithConfig(
		nil, client.DefaultTransportConfig().
			WithHost(env.AnalysisAPIHost).
			WithBasePath("/"+env.AnalysisAPIPath))

	alertsDB = &table.AlertsTable{
		AlertsTableName:             env.AlertsTableName,
		Client:                      dynamodb.New(awsSession),
		EventsTableName:             env.EventsTableName,
		RuleIDCreationTimeIndexName: env.RuleIndexName,
	}
}
