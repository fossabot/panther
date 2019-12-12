package pollers

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	resourcesapi "github.com/panther-labs/panther/api/resources/models"
	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	awspollers "github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/utils"
)

// Test Lambda Context

func testContext() context.Context {
	return lambdacontext.NewContext(
		context.Background(),
		&lambdacontext.LambdaContext{
			InvokedFunctionArn: "arn:aws:lambda:us-west-2:123456789123:function:snapshot-pollers:live",
			AwsRequestID:       "ad32d898-2a37-484d-9c50-3708c8fbc7d6",
		},
	)
}

// Mock Lambda Client

//type mockLambdaClient struct {
//	lambdaiface.LambdaAPI
//	mock.Mock
//}
//
//func (client *mockLambdaClient) Invoke(
//	input *lambda.InvokeInput,
//) (*lambda.InvokeOutput, error) {
//
//	args := client.Called(input)
//	return args.Get(0).(*lambda.InvokeOutput), args.Error(1)
//}

var (
	mockTime          = time.Time{}
	testIntegrationID = "0aab70c6-da66-4bb9-a83c-bbe8f5717fde"
)

func mockTimeFunc() time.Time {
	return mockTime
}

func TestBatchResources(t *testing.T) {
	var testResources []*resourcesapi.AddResourceEntry
	for i := 0; i < 1100; i++ {
		testResources = append(testResources, &resourcesapi.AddResourceEntry{
			Attributes:      &awsmodels.CloudTrailMeta{},
			ID:              "arn:aws:cloudtrail:region:account-id:trail/trailname",
			IntegrationID:   resourcesapi.IntegrationID(testIntegrationID),
			IntegrationType: resourcesapi.IntegrationTypeAws,
			Type:            "AWS.CloudTrail",
		})
	}

	testBatches := batchResources(testResources)
	require.NotEmpty(t, testBatches)
	assert.Len(t, testBatches, 3)
	assert.Len(t, testBatches[0], 500)
	assert.Len(t, testBatches[1], 500)
	assert.Len(t, testBatches[2], 100)
}

func TestHandlerNonExistentIntegration(t *testing.T) {
	t.Skip("skipping until resources-api mock is in place")
	testIntegrations := &pollermodels.ScanMsg{
		Entries: []*pollermodels.ScanEntry{
			{
				AWSAccountID:     aws.String("123456789012"),
				IntegrationID:    &testIntegrationID,
				ResourceID:       aws.String("arn:aws:s3:::test"),
				ResourceType:     aws.String("AWS.NonExistentResource.Type"),
				ScanAllResources: aws.Bool(false),
			},
		},
	}
	testIntegrationBytes, err := json.Marshal(testIntegrations)
	require.NoError(t, err)

	// mockLambda := &mockLambdaClient{}
	// lambdaClient = mockLambda

	sampleEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				AWSRegion:     "us-west-2",
				MessageId:     "702a0aba-ab1f-11e8-b09c-f218981400a1",
				ReceiptHandle: "AQEBCki01vLygW9L6Xq1hcSNR90swZdtgZHP1N5hEU1Dt22p66gQFxKEsVo7ObxpC+b/",
				Body:          string(testIntegrationBytes),
				Md5OfBody:     "d3673b20e6c009a81c73961b798f838a",
			},
		},
	}

	require.NoError(t, Handle(testContext(), sampleEvent))
}

func TestHandler(t *testing.T) {
	t.Skip("skipping until resources-api mock is in place")
	testIntegrations := &pollermodels.ScanMsg{
		Entries: []*pollermodels.ScanEntry{
			{
				AWSAccountID:     aws.String("123456789012"),
				IntegrationID:    &testIntegrationID,
				ScanAllResources: aws.Bool(true),
			},
		},
	}
	testIntegrationBytes, err := json.Marshal(testIntegrations)
	require.NoError(t, err)

	// Mock Lambda request/responses: Update Integration Start
	// updateIntegrationStartPayload := &apimodels.LambdaInput{
	// 	UpdateIntegrationLastScanStart: &apimodels.UpdateIntegrationLastScanStartInput{
	// 		IntegrationID:     &testIntegrationID,
	// 		LastScanStartTime: aws.Time(mockTime),
	// 		ScanStatus:        aws.String("scanning"),
	// 	},
	// }
	// updateIntegrationStartBytes, err := json.Marshal(updateIntegrationStartPayload)
	// require.NoError(t, err)
	// updateIntegrationStartInvokeInput := &lambda.InvokeInput{
	// 	FunctionName: &snapshotAPIFunctionName,
	// 	Payload:      updateIntegrationStartBytes,
	// }
	//
	// // Mock Lambda request/responses: Update Integration End
	// updateIntegrationEndPayload := &apimodels.LambdaInput{
	// 	UpdateIntegrationLastScanEnd: &apimodels.UpdateIntegrationLastScanEndInput{
	// 		IntegrationID:        &testIntegrationID,
	// 		LastScanEndTime:      aws.Time(mockTime),
	// 		ScanStatus:           aws.String("ok"),
	// 		LastScanErrorMessage: aws.String(""),
	// 	},
	// }
	// updateIntegrationEndBytes, err := json.Marshal(updateIntegrationEndPayload)
	// require.NoError(t, err)
	// updateIntegrationEndInvokeInput := &lambda.InvokeInput{
	// 	FunctionName: &snapshotAPIFunctionName,
	// 	Payload:      updateIntegrationEndBytes,
	// }
	//
	// mockLambda := &mockLambdaClient{}
	// genericInvokeOutSuccess := &lambda.InvokeOutput{StatusCode: aws.Int64(200)}
	// mockLambda.On("Invoke", updateIntegrationStartInvokeInput).Return(genericInvokeOutSuccess, nil)
	// mockLambda.On("Invoke", updateIntegrationEndInvokeInput).Return(genericInvokeOutSuccess, nil)
	// lambdaClient = mockLambda

	// Setup ACM client and function mocks
	awstest.MockAcmForSetup = awstest.BuildMockAcmSvcAll()

	// Setup CloudFormation client and function mocks
	awstest.MockCloudFormationForSetup = awstest.BuildMockCloudFormationSvcAll()

	// Setup CloudWatchLogs client and function mocks
	awstest.MockCloudWatchLogsForSetup = awstest.BuildMockCloudWatchLogsSvcAll()

	// Setup CloudTrail client and function mocks
	awstest.MockCloudTrailForSetup = awstest.BuildMockCloudTrailSvcAll()

	// Setup IAM client and function mocks
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcAll()

	// Setup Lambda client and function mocks
	awstest.MockLambdaForSetup = awstest.BuildMockLambdaSvcAll()

	// Setup S3 client and function mocks
	awstest.MockS3ForSetup = awstest.BuildMockS3SvcAll()

	// Setup EC2 client and function mocks
	awstest.MockEC2ForSetup = awstest.BuildMockEC2SvcAll()

	// Setup KMS client and function mocks
	awstest.MockKmsForSetup = awstest.BuildMockKmsSvcAll()

	// Setup ConfigService client with mock functions
	awstest.MockConfigServiceForSetup = awstest.BuildMockConfigServiceSvcAll()

	// Setup ELBV2 client and function mocks
	awstest.MockElbv2ForSetup = awstest.BuildMockElbv2SvcAll()

	// Setup WAF client and function mocks
	awstest.MockWafForSetup = awstest.BuildMockWafSvcAll()

	// Setup WAF Regional client and function mocks
	awstest.MockWafRegionalForSetup = awstest.BuildMockWafRegionalSvcAll()

	// Setup GuardDuty client and function mocks
	awstest.MockGuardDutyForSetup = awstest.BuildMockGuardDutySvcAll()

	// Setup DynamoDB client and function mocks
	awstest.MockDynamoDBForSetup = awstest.BuildMockDynamoDBSvcAll()

	// Setup DynamoDB client and function mocks
	awstest.MockApplicationAutoScalingForSetup = awstest.BuildMockApplicationAutoScalingSvcAll()

	// Setup RDS client and function mocks
	awstest.MockRdsForSetup = awstest.BuildMockRdsSvcAll()

	// Setup Redshift client and function mocks
	awstest.MockRedshiftForSetup = awstest.BuildMockRedshiftSvcAll()

	mockStsClient := &awstest.MockSTS{}
	mockStsClient.
		On("GetCallerIdentity", &sts.GetCallerIdentityInput{}).
		Return(
			&sts.GetCallerIdentityOutput{
				Account: aws.String("123456789012"),
				Arn:     aws.String("arn:aws:iam::123456789012:role/PantherAuditRole"),
				UserId:  aws.String("mockUserId"),
			},
			nil,
		)
	awstest.MockSTSForSetup = mockStsClient

	awspollers.AcmClientFunc = awstest.SetupMockAcm
	awspollers.ApplicationAutoScalingClientFunc = awstest.SetupMockApplicationAutoScaling
	awspollers.CloudTrailClientFunc = awstest.SetupMockCloudTrail
	awspollers.CloudWatchLogsClientFunc = awstest.SetupMockCloudWatchLogs
	awspollers.CloudFormationClientFunc = awstest.SetupMockCloudFormation
	awspollers.ConfigServiceClientFunc = awstest.SetupMockConfigService
	awspollers.DynamoDBClientFunc = awstest.SetupMockDynamoDB
	awspollers.EC2ClientFunc = awstest.SetupMockEC2
	awspollers.Elbv2ClientFunc = awstest.SetupMockElbv2
	awspollers.GuardDutyClientFunc = awstest.SetupMockGuardDuty
	awspollers.IAMClientFunc = awstest.SetupMockIAM
	awspollers.KmsClientFunc = awstest.SetupMockKms
	awspollers.LambdaClientFunc = awstest.SetupMockLambda
	awspollers.RDSClientFunc = awstest.SetupMockRds
	awspollers.RedshiftClientFunc = awstest.SetupMockRedshift
	awspollers.S3ClientFunc = awstest.SetupMockS3
	awspollers.WafClientFunc = awstest.SetupMockWaf
	awspollers.WafRegionalClientFunc = awstest.SetupMockWafRegional

	awspollers.AssumeRoleFunc = awstest.AssumeRoleMock
	awspollers.STSClientFunc = awstest.SetupMockSTSClient
	awspollers.AssumeRoleProviderFunc = awstest.STSAssumeRoleProviderMock

	// Time mock
	utils.TimeNowFunc = mockTimeFunc

	sampleEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				AWSRegion:     "us-west-2",
				MessageId:     "702a0aba-ab1f-11e8-b09c-f218981400a1",
				ReceiptHandle: "AQEBCki01vLygW9L6Xq1hcSNR90swZdtgZHP1N5hEU1Dt22p66gQFxKEsVo7ObxpC+b/",
				Body:          string(testIntegrationBytes),
				Md5OfBody:     "d3673b20e6c009a81c73961b798f838a",
			},
		},
	}

	require.NoError(t, Handle(testContext(), sampleEvent))
}
