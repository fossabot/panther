package processor

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

const (
	refreshInterval         = 2 * time.Minute
	snapshotAPIFunctionName = "panther-snapshot-api"
)

var (
	// Keyed on accountID
	accounts            = make(map[string]*models.SourceIntegration)
	accountsLastUpdated time.Time
	// Setup the clients to talk to the Snapshot API
	sess                               = session.Must(session.NewSession())
	lambdaClient lambdaiface.LambdaAPI = lambda.New(sess)
)

func resetAccountCache() {
	accounts = make(map[string]*models.SourceIntegration)
}

func refreshAccounts() error {
	if len(accounts) != 0 && accountsLastUpdated.Add(refreshInterval).After(time.Now()) {
		zap.L().Info("using cached accounts")
		return nil
	}

	zap.L().Info("populating account cache")
	input := &models.LambdaInput{
		ListIntegrations: &models.ListIntegrationsInput{
			IntegrationType: aws.String("aws-scan"),
		},
	}
	var output []*models.SourceIntegration
	err := genericapi.Invoke(lambdaClient, snapshotAPIFunctionName, input, &output)
	if err != nil {
		return err
	}

	for _, integration := range output {
		accounts[*integration.AWSAccountID] = integration
	}
	accountsLastUpdated = time.Now()

	return nil
}
