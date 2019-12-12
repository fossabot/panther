package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"

	apimodels "github.com/panther-labs/panther/api/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

const dynamoDBServiceNameSpace = "dynamodb"

// Set as variables to be overridden in testing
var (
	DynamoDBClientFunc               = setupDynamoDBClient
	ApplicationAutoScalingClientFunc = setupApplicationAutoScalingClient
)

func setupDynamoDBClient(sess *session.Session, cfg *aws.Config) interface{} {
	cfg.MaxRetries = aws.Int(MaxRetries)
	return dynamodb.New(sess, cfg)
}

func setupApplicationAutoScalingClient(sess *session.Session, cfg *aws.Config) interface{} {
	return applicationautoscaling.New(sess, cfg)
}

// PollDynamoDBTable polls a single DynamoDB Table resource
func PollDynamoDBTable(
	pollerResourceInput *awsmodels.ResourcePollerInput,
	resourceARN arn.ARN,
	_ *pollermodels.ScanEntry,
) interface{} {

	dynamoClient := getClient(pollerResourceInput, "dynamodb", resourceARN.Region).(dynamodbiface.DynamoDBAPI)
	autoscalingClient := getClient(
		pollerResourceInput, "applicationautoscaling", resourceARN.Region).(applicationautoscalingiface.ApplicationAutoScalingAPI)
	table := strings.Replace(resourceARN.Resource, "table/", "", 1)

	snapshot := buildDynamoDBTableSnapshot(dynamoClient, autoscalingClient, aws.String(table))
	if snapshot == nil {
		return nil
	}
	snapshot.Region = aws.String(resourceARN.Region)
	snapshot.AccountID = aws.String(resourceARN.AccountID)
	return snapshot
}

// listTables returns a list of all Dynamo DB tables in the account
func listTables(dynamoDBSvc dynamodbiface.DynamoDBAPI) (tables []*string) {
	err := dynamoDBSvc.ListTablesPages(&dynamodb.ListTablesInput{},
		func(page *dynamodb.ListTablesOutput, lastPage bool) bool {
			tables = append(tables, page.TableNames...)
			return true
		})
	if err != nil {
		utils.LogAWSError("DynamoDB.ListTablesPages", err)
	}
	return
}

// describeTable provides detailed information about a given DynamoDB table
func describeTable(dynamoDBSvc dynamodbiface.DynamoDBAPI, name *string) *dynamodb.TableDescription {
	out, err := dynamoDBSvc.DescribeTable(&dynamodb.DescribeTableInput{TableName: name})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				zap.L().Warn("tried to scan non-existent resource",
					zap.String("resource", *name),
					zap.String("resourceType", awsmodels.DynamoDBTableSchema))
				return nil
			}
		}
		utils.LogAWSError("DynamoDB.DescribeTable", err)
		return nil
	}

	return out.Table
}

// describeTimeToLive provides time to live configuration information
func describeTimeToLive(dynamoDBSvc dynamodbiface.DynamoDBAPI, name *string) (*dynamodb.TimeToLiveDescription, error) {
	out, err := dynamoDBSvc.DescribeTimeToLive(&dynamodb.DescribeTimeToLiveInput{TableName: name})
	if err != nil {
		utils.LogAWSError("DynamoDB.DescribeTimeToLive", err)
		return nil, err
	}

	return out.TimeToLiveDescription, nil
}

// listTagsOfResource returns the tags for a given DynamoDB table
func listTagsOfResource(dynamoDBSvc dynamodbiface.DynamoDBAPI, arn *string) ([]*dynamodb.Tag, error) {
	out, err := dynamoDBSvc.ListTagsOfResource(&dynamodb.ListTagsOfResourceInput{ResourceArn: arn})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "AccessDeniedException" {
				zap.L().Error(
					"AccessDeniedException, additional permissions were not granted",
					zap.String("API", "DynamoDB.ListTagsOfResource"))
				return nil, err
			}
		}

		utils.LogAWSError("DynamoDB.ListTagsOfResource ", err)
		return nil, err
	}

	return out.Tags, nil
}

// describeScalableTargets provides information about autoscaling for a given resource
// Gathers autoscaling configuration on both a DynamoDB table and its Global Secondary Indices (GSI's)
func describeScalableTargets(
	applicationAutoScalingSvc applicationautoscalingiface.ApplicationAutoScalingAPI,
	resourceIDs []*string,
) (autoscaling []*applicationautoscaling.ScalableTarget, err error) {

	input := &applicationautoscaling.DescribeScalableTargetsInput{
		ResourceIds:      resourceIDs,
		ServiceNamespace: aws.String(dynamoDBServiceNameSpace),
	}
	err = applicationAutoScalingSvc.DescribeScalableTargetsPages(input,
		func(page *applicationautoscaling.DescribeScalableTargetsOutput, lastPage bool) bool {
			autoscaling = append(autoscaling, page.ScalableTargets...)
			return true
		})
	if err != nil {
		utils.LogAWSError("ApplicationAutoScaling.DescribeScalableTargetsPages", err)
		return
	}

	return
}

// buildDynamoDBTableSnapshot builds a snapshot of a DynamoDB table, including information about its
// Global Secondary Indices (GSI's) and any applicable autoscaling information for the table and GSI's
func buildDynamoDBTableSnapshot(
	dynamoDBSvc dynamodbiface.DynamoDBAPI,
	applicationAutoScalingSvc applicationautoscalingiface.ApplicationAutoScalingAPI,
	tableName *string,
) *awsmodels.DynamoDBTable {

	description := describeTable(dynamoDBSvc, tableName)
	// Some type of error occurred, it's already been logged appropriately in describeTable but we
	// cannot continue building this resource.
	if description == nil {
		return nil
	}

	table := &awsmodels.DynamoDBTable{
		GenericResource: awsmodels.GenericResource{
			ResourceType: aws.String(awsmodels.DynamoDBTableSchema),
			ResourceID:   description.TableArn,
			TimeCreated:  utils.DateTimeFormat(*description.CreationDateTime),
		},
		GenericAWSResource: awsmodels.GenericAWSResource{
			Name: tableName,
			ARN:  description.TableArn,
			ID:   description.TableId,
		},
		AttributeDefinitions:   description.AttributeDefinitions,
		BillingModeSummary:     description.BillingModeSummary,
		GlobalSecondaryIndexes: description.GlobalSecondaryIndexes,
		ItemCount:              description.ItemCount,
		KeySchema:              description.KeySchema,
		LatestStreamArn:        description.LatestStreamArn,
		LatestStreamLabel:      description.LatestStreamLabel,
		LocalSecondaryIndexes:  description.LocalSecondaryIndexes,
		ProvisionedThroughput:  description.ProvisionedThroughput,
		RestoreSummary:         description.RestoreSummary,
		SSEDescription:         description.SSEDescription,
		StreamSpecification:    description.StreamSpecification,
		TableSizeBytes:         description.TableSizeBytes,
		TableStatus:            description.TableStatus,
	}

	tableID := aws.String("table/" + *tableName)
	resourceIDs := []*string{tableID}

	ttl, err := describeTimeToLive(dynamoDBSvc, tableName)
	if err == nil {
		table.TimeToLiveDescription = ttl
	}

	for _, index := range description.GlobalSecondaryIndexes {
		indexID := aws.String(*tableID + "/index/" + *index.IndexName)
		resourceIDs = append(resourceIDs, indexID)
	}

	tags, err := listTagsOfResource(dynamoDBSvc, table.ARN)
	if err == nil {
		table.Tags = utils.ParseTagSlice(tags)
	}

	autoScalingDescriptions, err := describeScalableTargets(applicationAutoScalingSvc, resourceIDs)
	if err != nil {
		utils.LogAWSError("ApplicationAutoScaling.DescribeScalableTargets", err)
		return table
	}
	table.AutoScalingDescriptions = autoScalingDescriptions

	return table
}

// PollDynamoDBTables gathers information on each Dynamo DB Table for an AWS account.
func PollDynamoDBTables(pollerInput *awsmodels.ResourcePollerInput) ([]*apimodels.AddResourceEntry, error) {
	zap.L().Debug("starting DynamoDB Table resource poller")
	dynamoDBTableSnapshots := make(map[string]*awsmodels.DynamoDBTable)

	for _, regionID := range utils.GetServiceRegions(pollerInput.Regions, "dynamodb") {
		sess := session.Must(session.NewSession(&aws.Config{Region: regionID}))
		creds, err := AssumeRoleFunc(pollerInput, sess)
		if err != nil {
			return nil, err
		}

		config := &aws.Config{Credentials: creds}
		dynamoDBSvc := DynamoDBClientFunc(sess, config).(dynamodbiface.DynamoDBAPI)
		applicationAutoScalingSvc := ApplicationAutoScalingClientFunc(sess, config).(applicationautoscalingiface.ApplicationAutoScalingAPI)

		// Start with generating a list of all tables
		tables := listTables(dynamoDBSvc)
		if len(tables) == 0 {
			zap.L().Debug("no DynamoDB tables found.", zap.String("region", *regionID))
			continue
		}

		for _, table := range tables {
			dynamoDBTable := buildDynamoDBTableSnapshot(
				dynamoDBSvc,
				applicationAutoScalingSvc,
				table,
			)
			if dynamoDBTable == nil {
				continue
			}
			dynamoDBTable.AccountID = aws.String(pollerInput.AuthSourceParsedARN.AccountID)
			dynamoDBTable.Region = regionID

			if _, ok := dynamoDBTableSnapshots[*dynamoDBTable.ARN]; !ok {
				dynamoDBTableSnapshots[*dynamoDBTable.ARN] = dynamoDBTable
			} else {
				zap.L().Info(
					"overwriting existing DynamoDB Table snapshot",
					zap.String("resourceID", *dynamoDBTable.ARN),
				)
				dynamoDBTableSnapshots[*dynamoDBTable.ARN] = dynamoDBTable
			}
		}
	}

	resources := make([]*apimodels.AddResourceEntry, 0, len(dynamoDBTableSnapshots))
	for resourceID, table := range dynamoDBTableSnapshots {
		resources = append(resources, &apimodels.AddResourceEntry{
			Attributes:      table,
			ID:              apimodels.ResourceID(resourceID),
			IntegrationID:   apimodels.IntegrationID(*pollerInput.IntegrationID),
			IntegrationType: apimodels.IntegrationTypeAws,
			Type:            awsmodels.DynamoDBTableSchema,
		})
	}

	return resources, nil
}
