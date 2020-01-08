package processor

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

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyKMS(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awskeymanagementservice.html
	if strings.HasPrefix(eventName, "Decrypt") ||
		strings.HasPrefix(eventName, "GenerateDataKey") ||
		strings.HasPrefix(eventName, "Encrypt") {

		zap.L().Debug("kms: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	var keyARN string
	switch eventName {
	/*
		Missing (not sure if needed in all cases):
			(Connect/Create/Delete/Update)CustomKeyStore
			(Delete/Import)KeyMaterial
			(Retire/Revoke)Grant
	*/
	case "CancelKeyDeletion", "CreateGrant", "CreateKey", "DisableKey", "DisableKeyRotation", "EnableKey",
		"EnableKeyRotation", "PutKeyPolicy", "ScheduleKeyDeletion", "TagResource", "UntagResource",
		"UpdateAlias", "UpdateKeyDescription":
		keyARN = detail.Get("resources").Array()[0].Get("ARN").Str
	case "CreateAlias", "DeleteAlias":
		resources := detail.Get("resources").Array()
		for _, resource := range resources {
			resourceARN, err := arn.Parse(resource.Get("ARN").Str)
			if err != nil {
				zap.L().Error("kms: unable to extract ARN", zap.String("eventName", eventName))
				return nil
			}
			if strings.HasPrefix(resourceARN.Resource, "key/") {
				keyARN = resourceARN.String()
			}
		}
	default:
		zap.L().Warn("kms: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteKey",
		EventName:    eventName,
		ResourceID:   keyARN,
		ResourceType: schemas.KmsKeySchema,
	}}
}
