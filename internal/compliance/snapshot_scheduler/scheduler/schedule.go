package scheduler

import (
	"time"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"go.uber.org/zap"

	snapshotapi "github.com/panther-labs/panther/internal/compliance/snapshot_api/api"
	"github.com/panther-labs/panther/pkg/genericapi"
)

const snapshotAPIFunctionName = "panther-snapshot-api:live"

var (
	sess                               = session.Must(session.NewSession())
	lambdaClient lambdaiface.LambdaAPI = lambda.New(sess)
)

// PollAndIssueNewScans sends messages to the snapshot-pollers when new scans need to start.
func PollAndIssueNewScans() error {
	enabledIntegrations, err := getEnabledIntegrations()
	if err != nil {
		return err
	}
	if len(enabledIntegrations) == 0 {
		zap.L().Info("no scans to schedule")
		return nil
	}

	zap.L().Info("loaded enabled integrations", zap.Int("count", len(enabledIntegrations)))
	var integrationsToScan []*models.SourceIntegrationMetadata

	for _, integration := range enabledIntegrations {
		// Only add new scans if needed
		if (scanIntervalElapsed(integration) && scanIsNotOngoing(integration)) || scanIsStuck(integration) {
			integrationsToScan = append(integrationsToScan, integration.SourceIntegrationMetadata)
		} else {
			zap.L().Debug("skipping integration", zap.String("integrationID", *integration.IntegrationID))
		}
	}

	return snapshotapi.ScanAllResources(integrationsToScan)
}

// getEnabledIntegrations lists enabled integrations from the snapshot-api.
func getEnabledIntegrations() (integrations []*models.SourceIntegration, err error) {
	err = genericapi.Invoke(
		lambdaClient,
		snapshotAPIFunctionName,
		&models.LambdaInput{ListIntegrations: &models.ListIntegrationsInput{
			IntegrationType: aws.String("aws-scan"),
		}},
		&integrations,
	)
	if err != nil {
		return
	}

	return
}

// scanIsStuck checks if an integration's is stuck in the "scanning" state.
func scanIsStuck(integration *models.SourceIntegration) bool {
	// Accounts for a new integration that has not completed a scan
	if integration.SourceIntegrationStatus == nil || integration.LastScanEndTime == nil {
		return false
	}

	return *integration.ScanStatus == models.StatusScanning && scanIntervalElapsed(integration)
}

// scanIsNotOngoing checks if an integration's snapshot is currently running.
func scanIsNotOngoing(integration *models.SourceIntegration) bool {
	if integration.SourceIntegrationStatus == nil {
		return true
	}

	return *integration.ScanStatus != models.StatusScanning
}

// scanIntervalElapsed determines if a new scan needs to be started based on the configured interval.
func scanIntervalElapsed(integration *models.SourceIntegration) bool {
	// Account for cases when a scan has never ran.
	if integration.SourceIntegrationScanInformation == nil {
		return true
	}

	if integration.SourceIntegrationScanInformation.LastScanEndTime == nil {
		return true
	}

	intervalMins := time.Duration(*integration.SourceIntegrationMetadata.ScanIntervalMins) * time.Minute
	return time.Since(*integration.LastScanEndTime) >= intervalMins
}
