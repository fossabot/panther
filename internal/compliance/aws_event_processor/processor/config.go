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

	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyConfig(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// We need to add more config resources, just a config recorder is too high level
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awsconfig.html
	if eventName == "PutAggregationAuthorization" ||
		eventName == "PutConfigurationAggregator" ||
		eventName == "PutDeliveryChannel" ||
		eventName == "PutEvaluations" ||
		eventName == "PutRemediationConfigurations" ||
		eventName == "PutRetentionConfiguration" ||
		eventName == "StartRemediationExecution" ||
		eventName == "TagResource" ||
		eventName == "UntagResource" ||
		eventName == "DeleteDeliveryChannel" ||
		eventName == "DeleteEvaluationResults" ||
		eventName == "DeletePendingAggregationRequest" ||
		eventName == "DeleteRemediationConfiguration" ||
		eventName == "DeleteRetentionConfiguration" ||
		eventName == "DeliverConfigSnapshot" ||
		eventName == "DeleteAggregationAuthorization" ||
		eventName == "DeleteConfigRule" ||
		eventName == "DeleteConfigurationAggregator" ||
		eventName == "PutConfigRule" {

		zap.L().Debug("config: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	switch eventName {
	case "StartConfigRulesEvaluation", "StartConfigurationRecorder", "StopConfigurationRecorder":
		// This case handles when a recorder is updated in a way that does not require a full account
		// scan to update the config meta resource
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			ResourceID: strings.Join([]string{
				accountID,
				detail.Get("awsRegion").Str,
				schemas.ConfigServiceSchema,
			}, ":"),
			ResourceType: schemas.ConfigServiceSchema,
		}}
	case "PutConfigurationRecorder":
		// This case handles when a recorder is updated in a way that requires a full account scan
		// in order to update the config meta resource
		return []*resourceChange{{
			AwsAccountID: accountID,
			EventName:    eventName,
			ResourceType: schemas.ConfigServiceSchema,
		}}
	case "DeleteConfigurationRecorder":
		// Special case where need to queue both a delete action and a meta re-scan
		return []*resourceChange{
			{
				AwsAccountID: accountID,
				Delete:       true,
				EventName:    eventName,
				ResourceID: strings.Join([]string{
					accountID,
					detail.Get("awsRegion").Str,
					schemas.ConfigServiceSchema,
				}, ":"),
				ResourceType: schemas.ConfigServiceSchema,
			},
			{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceType: schemas.ConfigServiceSchema,
			}}
	default:
		zap.L().Warn("config: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}
}
