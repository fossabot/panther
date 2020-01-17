package api

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
