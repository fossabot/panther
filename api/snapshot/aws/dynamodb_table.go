package aws

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
