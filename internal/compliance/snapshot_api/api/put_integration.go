package api

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	awspoller "github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws"
	"github.com/panther-labs/panther/pkg/awsbatch/sqsbatch"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// PutIntegration adds a set of new integrations in a batch.
func (API) PutIntegration(input *models.PutIntegrationInput) ([]*models.SourceIntegrationMetadata, error) {
	newIntegrations := make([]*models.SourceIntegrationMetadata, len(input.Integrations))

	// Generate the new integrations
	for i, integration := range input.Integrations {
		newIntegrations[i] = generateNewIntegration(integration)
	}

	// Batch write to DynamoDB
	if err := db.BatchPutSourceIntegrations(newIntegrations); err != nil {
		return nil, err
	}

	// Return early to skip sending to the snapshot queue
	if aws.BoolValue(input.SkipScanQueue) {
		return newIntegrations, nil
	}

	var integrationsToScan []*models.SourceIntegrationMetadata
	for _, integration := range newIntegrations {
		//We don't want to trigger scanning for aws-s3 type integrations
		if aws.StringValue(integration.IntegrationType) == models.IntegrationTypeAWS3 {
			continue
		}
		integrationsToScan = append(integrationsToScan, integration)
	}

	// Add to the Snapshot queue
	return newIntegrations, ScanAllResources(integrationsToScan)
}

// ScanAllResources schedules scans for each resource type for each integration.
//
// Each resource type is sent within its own SQS message.
func ScanAllResources(integrations []*models.SourceIntegrationMetadata) error {
	var sqsEntries []*sqs.SendMessageBatchRequestEntry

	// For each integration, add a ScanMsg to the queue per service
	for _, integration := range integrations {
		if !*integration.ScanEnabled {
			continue
		}

		for resourceType := range awspoller.ServicePollers {
			scanMsg := &pollermodels.ScanMsg{
				Entries: []*pollermodels.ScanEntry{
					{
						AWSAccountID:  integration.AWSAccountID,
						IntegrationID: integration.IntegrationID,
						ResourceType:  aws.String(resourceType),
					},
				},
			}

			messageBodyBytes, err := jsoniter.MarshalToString(scanMsg)
			if err != nil {
				return &genericapi.InternalError{Message: err.Error()}
			}

			sqsEntries = append(sqsEntries, &sqs.SendMessageBatchRequestEntry{
				// Generates an ID of: IntegrationID-AWSResourceType
				Id: aws.String(
					*integration.IntegrationID + "-" + strings.Replace(resourceType, ".", "", -1),
				),
				MessageBody: aws.String(messageBodyBytes),
			})
		}
	}

	zap.L().Info(
		"scheduling new scans",
		zap.String("queueUrl", snapshotPollersQueueURL),
		zap.Int("count", len(sqsEntries)),
	)

	// Batch send all the messages to SQS
	return sqsbatch.SendMessageBatch(SQSClient, maxElapsedTime, &sqs.SendMessageBatchInput{
		Entries:  sqsEntries,
		QueueUrl: &snapshotPollersQueueURL,
	})
}

func generateNewIntegration(input *models.PutIntegrationSettings) *models.SourceIntegrationMetadata {
	return &models.SourceIntegrationMetadata{
		AWSAccountID:     input.AWSAccountID,
		CreatedAtTime:    aws.Time(time.Now()),
		CreatedBy:        input.UserID,
		IntegrationID:    aws.String(uuid.New().String()),
		IntegrationLabel: input.IntegrationLabel,
		IntegrationType:  input.IntegrationType,
		ScanEnabled:      input.ScanEnabled,
		ScanIntervalMins: input.ScanIntervalMins,
		// For log analysis integrations
		SourceSnsTopicArn:    input.SourceSnsTopicArn,
		LogProcessingRoleArn: input.LogProcessingRoleArn,
	}
}
