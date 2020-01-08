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

func classifyDynamoDB(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_amazondynamodb.html
	if eventName == "BatchGetItem" ||
		eventName == "ConditionCheckItem" ||
		eventName == "DeleteBackup" ||
		eventName == "DeleteItem" ||
		eventName == "PutItem" ||
		eventName == "Query" ||
		eventName == "Scan" ||
		eventName == "UpdateItem" ||
		eventName == "BatchWriteItem" {

		zap.L().Debug("dynamodb: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	dynamoARN := arn.ARN{
		Partition: "aws",
		Service:   "dynamodb",
		Region:    detail.Get("awsRegion").Str,
		AccountID: accountID,
		Resource:  "table/",
	}
	switch eventName {
	case "CreateBackup", "CreateTable", "DeleteTable", "UpdateContinuousBackups", "UpdateTable", "UpdateTimeToLive":
		dynamoARN.Resource += detail.Get("requestParameters.tableName").Str
	case "CreateGlobalTable":
		var tables []*resourceChange
		dynamoARN.Resource += detail.Get("requestParameters.globalTableName").Str
		for _, repl := range detail.Get("requestParameters.replicationGroup").Array() {
			dynamoARN.Region = repl.Get("regionName").Str
			tables = append(tables, &resourceChange{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceID:   dynamoARN.String(),
				ResourceType: schemas.DynamoDBTableSchema,
			})
		}
		return tables
	case "UpdateGlobalTable":
		// What an odd API call...
		var tables []*resourceChange
		dynamoARN.Resource += detail.Get("requestParameters.globalTableName").Str
		for _, update := range detail.Get("requestParameters.replicaUpdates").Array() {
			region := update.Get("create.regionName").Str
			if region != "" {
				dynamoARN.Region = region
				tables = append(tables, &resourceChange{
					AwsAccountID: accountID,
					EventName:    eventName,
					ResourceID:   dynamoARN.String(),
					ResourceType: schemas.DynamoDBTableSchema,
				})
			}
			region = update.Get("delete.regionName").Str
			if region != "" {
				dynamoARN.Region = region
				tables = append(tables, &resourceChange{
					AwsAccountID: accountID,
					EventName:    eventName,
					ResourceID:   dynamoARN.String(),
					ResourceType: schemas.DynamoDBTableSchema,
				})
			}
		}
		return tables
	case "UpdateGlobalTableSettings":
		// Untested, feels right though
		var tables []*resourceChange
		dynamoARN.Resource += detail.Get("requestParameters.globalTableName").Str
		for _, replica := range detail.Get("responseElements.replicaSettings").Array() {
			dynamoARN.Region = replica.Get("regionName").Str
			tables = append(tables, &resourceChange{
				AwsAccountID: accountID,
				EventName:    eventName,
				ResourceID:   dynamoARN.String(),
				ResourceType: schemas.DynamoDBTableSchema,
			})
		}
		return tables
	case "RestoreTableFromBackup", "RestoreTableToPointInTime":
		dynamoARN.Resource += detail.Get("requestParameters.targetTableName").Str
	case "TagResource", "UntagResource":
		tableARN, err := arn.Parse(detail.Get("requestParameters.resourceArn").Str)
		if err != nil {
			zap.L().Error("dynamodb: error parsing ARN", zap.Error(err))
			return nil
		}
		dynamoARN = tableARN
	default:
		zap.L().Warn("dynamodb: encountered unknown event name", zap.String("eventName", eventName))
		return nil
	}

	return []*resourceChange{{
		AwsAccountID: accountID,
		Delete:       eventName == "DeleteTable",
		EventName:    eventName,
		ResourceID:   dynamoARN.String(),
		ResourceType: schemas.DynamoDBTableSchema,
	}}
}
