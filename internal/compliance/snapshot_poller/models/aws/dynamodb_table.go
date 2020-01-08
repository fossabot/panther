package aws

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
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	DynamoDBTableSchema = "AWS.DynamoDB.Table"
)

// DynamoDBTable contains all the information about a Dynamo DB table
type DynamoDBTable struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from dynamodb.TableDescription
	AttributeDefinitions   []*dynamodb.AttributeDefinition
	BillingModeSummary     *dynamodb.BillingModeSummary
	GlobalSecondaryIndexes []*dynamodb.GlobalSecondaryIndexDescription
	ItemCount              *int64
	KeySchema              []*dynamodb.KeySchemaElement
	LatestStreamArn        *string
	LatestStreamLabel      *string
	LocalSecondaryIndexes  []*dynamodb.LocalSecondaryIndexDescription
	ProvisionedThroughput  *dynamodb.ProvisionedThroughputDescription
	RestoreSummary         *dynamodb.RestoreSummary
	SSEDescription         *dynamodb.SSEDescription
	StreamSpecification    *dynamodb.StreamSpecification
	TableSizeBytes         *int64
	TableStatus            *string

	// Additional fields
	//
	// Both a Dynamo Table and its Global Secondary Indices can be an auto scaling target
	// This is a list of a table and its indices autoscaling configurations (if they exist)
	//
	AutoScalingDescriptions []*applicationautoscaling.ScalableTarget
	TimeToLiveDescription   *dynamodb.TimeToLiveDescription
}
