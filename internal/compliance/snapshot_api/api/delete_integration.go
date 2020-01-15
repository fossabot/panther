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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/pkg/genericapi"
)

// DeleteIntegration deletes a specific integration.
func (API) DeleteIntegration(input *models.DeleteIntegrationInput) (err error) {
	var integrationForDeletePermissions *models.SourceIntegrationMetadata
	defer func() {
		if err != nil && integrationForDeletePermissions != nil {
			// In case we have already removed the Permissions from SQS but some other operation failed
			// re-add the permissions
			if undoErr := updateLogProcessingPermissions(integrationForDeletePermissions); undoErr != nil {
				zap.L().Error("failed to re-add SQS permission for integration. SQS is missing permissions that have to be added manually",
					zap.String("integrationId", *integrationForDeletePermissions.IntegrationID),
					zap.Error(undoErr),
					zap.Error(err))
			}
		}
	}()

	var integration *models.SourceIntegrationMetadata
	integration, err = db.GetIntegration(input.IntegrationID)
	if err != nil {
		zap.L().Warn("failed to get integration", zap.String("integrationId", *input.IntegrationID), zap.Error(err))
		return err
	}

	if integration == nil {
		return &genericapi.DoesNotExistError{Message: "Integration does not exist"}
	}

	if *integration.IntegrationType == models.IntegrationTypeAWS3 {
		if err = removeLogProcessingPermissions(integration); err != nil {
			zap.L().Warn("failed to remove permission from SQS queue for integration",
				zap.String("integrationId", *input.IntegrationID),
				zap.Error(err))
			return err
		}
		integrationForDeletePermissions = integration
	}
	err = db.DeleteIntegrationItem(input)
	return err
}

func removeLogProcessingPermissions(input *models.SourceIntegrationMetadata) error {
	removePermissionInput := &sqs.RemovePermissionInput{
		QueueUrl: aws.String(logAnalysisQueueURL),
		Label:    input.IntegrationID,
	}
	_, err := SQSClient.RemovePermission(removePermissionInput)
	return err
}
