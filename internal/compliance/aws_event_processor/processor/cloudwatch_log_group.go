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
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyCloudWatchLogGroup(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_amazoncloudwatchlogs.html
	if eventName == "CancelExportTask" ||
		eventName == "CreateExportTask" ||
		eventName == "PutDestination" ||
		eventName == "PutDestinationPolicy" ||
		eventName == "PutLogEvents" ||
		eventName == "PutResourcePolicy" ||
		eventName == "StartQuery" ||
		eventName == "StopQuery" ||
		eventName == "TestMetricFilter" ||
		eventName == "CreateLogStream" ||
		eventName == "FilterLogEvents" {

		zap.L().Debug("loggroup: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	region := detail.Get("awsRegion").Str
	logGroupARN := arn.ARN{
		Partition: "aws",
		Service:   "logs",
		Region:    region,
		AccountID: accountID,
		Resource:  "log-group:",
	}
	switch eventName {
	case "AssociateKmsKey", "CreateLogGroup", "DeleteLogGroup", "DeleteLogStream", "DeleteMetricFilter",
		"DeleteRetentionPolicy", "DeleteSubscriptionFilter", "DisassociateKmsKey", "PutMetricFilter",
		"PutRetentionPolicy", "PutSubscriptionFilter", "TagLogGroup", "UntagLogGroup":
		// Not technically the correct resourceID, see classifyCloudFormation for a more detailed
		// explanation.
		logGroupARN.Resource += detail.Get("requestParameters.logGroupName").Str
	default:
		zap.L().Warn("loggroup: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteLogGroup",
		EventName:    eventName,
		ResourceID:   logGroupARN.String(),
		ResourceType: schemas.CloudWatchLogGroupSchema,
	}}
}
