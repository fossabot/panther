package api

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

const (
	sqsReceiveMessageAction = "ReceiveMessage"
)

// PutIntegration adds a set of new integrations in a batch.
func (API) PutIntegration(input *models.PutIntegrationInput) ([]*models.SourceIntegrationMetadata, error) {
	permissionsAddedForIntegrations := []*models.SourceIntegrationMetadata{}
	var err error
	defer func() {
		if err != nil {
			// In case there has been any error, try to undo granting of permissions to SQS queue.
			for _, integration := range permissionsAddedForIntegrations {
				if undoErr := removeLogProcessingPermissions(integration); undoErr != nil {
					zap.L().Error("failed to remove SQS permission for integration. SQS queue has additional permissions that have to be removed manually",
						zap.String("sqsPermissionLabel", *integration.IntegrationID),
						zap.Error(undoErr),
						zap.Error(err))
				}
			}
		}
	}()
	newIntegrations := make([]*models.SourceIntegrationMetadata, len(input.Integrations))

	// Generate the new integrations
	for i, integration := range input.Integrations {
		newIntegrations[i] = generateNewIntegration(integration)
	}

	for _, integration := range newIntegrations {
		err = updateLogProcessingPermissions(integration)
		if err != nil {
			return nil, err
		}
		permissionsAddedForIntegrations = append(permissionsAddedForIntegrations, integration)
	}

	// Batch write to DynamoDB
	if err = db.BatchPutSourceIntegrations(newIntegrations); err != nil {
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
	err = ScanAllResources(integrationsToScan)
	return newIntegrations, err
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
		S3Buckets: input.S3Buckets,
		KmsKeys:   input.KmsKeys,
	}
}

// updateLogProcessingPermissions updates Log Processor SQS queue to allow new account
// to send data to it.
func updateLogProcessingPermissions(input *models.SourceIntegrationMetadata) error {
	if *input.IntegrationType != models.IntegrationTypeAWS3 {
		return nil
	}
	permissionInput := &sqs.AddPermissionInput{
		AWSAccountIds: []*string{input.AWSAccountID},
		Actions:       aws.StringSlice([]string{sqsReceiveMessageAction}),
		QueueUrl:      aws.String(logAnalysisQueueURL),
		Label:         input.IntegrationID,
	}
	_, err := SQSClient.AddPermission(permissionInput)
	return err
}
