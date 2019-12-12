package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	snapshotapimodels "github.com/panther-labs/panther/api/snapshot"
)

//
// Mocks
//

// mockLambdaClient mocks the API calls to the snapshot-api.
type mockLambdaClient struct {
	lambdaiface.LambdaAPI
	mock.Mock
}

// Invoke is a mock method to invoke a Lambda function.
func (client *mockLambdaClient) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	args := client.Called(input)
	return args.Get(0).(*lambda.InvokeOutput), args.Error(1)
}

//
// Helpers
//

// getTestInvokeInput returns an example Lambda.Invoke input for the SnapshotAPI.
func getTestInvokeInput() *lambda.InvokeInput {
	input := &snapshotapimodels.LambdaInput{
		ListIntegrations: &snapshotapimodels.ListIntegrationsInput{
			IntegrationType: aws.String("aws-scan"),
		},
	}
	payload, err := jsoniter.Marshal(input)
	if err != nil {
		panic(err)
	}

	return &lambda.InvokeInput{
		FunctionName: aws.String("panther-snapshot-api:live"),
		Payload:      payload,
	}
}

// getTestInvokeOutput returns an example Lambda.Invoke response from the SnapshotAPI.
func getTestInvokeOutput(payload interface{}, statusCode int64) *lambda.InvokeOutput {
	payloadBytes, err := jsoniter.Marshal(payload)
	if err != nil {
		panic(err)
	}

	return &lambda.InvokeOutput{
		Payload:    payloadBytes,
		StatusCode: aws.Int64(statusCode),
	}
}

//
// Unit Tests
//

var (
	exampleIntegrations = []*snapshotapimodels.SourceIntegration{
		{
			SourceIntegrationMetadata: &snapshotapimodels.SourceIntegrationMetadata{
				IntegrationID:    aws.String("45c378a7-2e36-4b12-8e16-2d3c49ff1371"),
				IntegrationLabel: aws.String("ProdAWS"),
				IntegrationType:  aws.String("aws-scan"),
				ScanIntervalMins: aws.Int(60),
				ScanEnabled:      aws.Bool(true),
			},
			SourceIntegrationStatus: &snapshotapimodels.SourceIntegrationStatus{
				ScanStatus: aws.String("ok"),
			},
			SourceIntegrationScanInformation: &snapshotapimodels.SourceIntegrationScanInformation{
				LastScanEndTime:   aws.Time(time.Now().Add(time.Duration(-15) * time.Minute)),
				LastScanStartTime: aws.Time(time.Now().Add(time.Duration(-20) * time.Minute)),
			},
		},
		{
			SourceIntegrationMetadata: &snapshotapimodels.SourceIntegrationMetadata{
				IntegrationID:    aws.String("ebb4d69f-177b-4eff-a7a6-9251fdc72d21"),
				IntegrationLabel: aws.String("TestAWS"),
				IntegrationType:  aws.String("aws-scan"),
				ScanIntervalMins: aws.Int(30),
				ScanEnabled:      aws.Bool(true),
			},
			SourceIntegrationStatus: &snapshotapimodels.SourceIntegrationStatus{
				ScanStatus: aws.String("ok"),
			},
			SourceIntegrationScanInformation: &snapshotapimodels.SourceIntegrationScanInformation{
				LastScanEndTime:   aws.Time(time.Now().Add(time.Duration(-35) * time.Minute)),
				LastScanStartTime: aws.Time(time.Now().Add(time.Duration(-40) * time.Minute)),
			},
		},
		// A new integration that was recently added that has never been scanned.
		{
			SourceIntegrationMetadata: &snapshotapimodels.SourceIntegrationMetadata{
				IntegrationID:    aws.String("ebb4d69f-177b-4eff-a7a6-9251fdc72d21"),
				IntegrationLabel: aws.String("TestAWS"),
				IntegrationType:  aws.String("aws-scan"),
				ScanIntervalMins: aws.Int(30),
				ScanEnabled:      aws.Bool(true),
			},
		},
		// An integration with a scan in progress, started 20 minutes ago.
		{
			SourceIntegrationMetadata: &snapshotapimodels.SourceIntegrationMetadata{
				IntegrationID:    aws.String("9a171500-7794-4aaa-8b4a-19ce8e9ba4fb"),
				IntegrationLabel: aws.String("Staging AWS Account"),
				IntegrationType:  aws.String("aws-scan"),
				ScanIntervalMins: aws.Int(60),
				ScanEnabled:      aws.Bool(true),
			},
			SourceIntegrationStatus: &snapshotapimodels.SourceIntegrationStatus{
				ScanStatus: aws.String("scanning"),
			},
			SourceIntegrationScanInformation: &snapshotapimodels.SourceIntegrationScanInformation{
				LastScanStartTime: aws.Time(time.Now().Add(time.Duration(-20) * time.Minute)),
			},
		},
		// An integration with a scan in progress, stuck.
		{
			SourceIntegrationMetadata: &snapshotapimodels.SourceIntegrationMetadata{
				IntegrationID:    aws.String("2654cf7a-a13a-4b9b-8b4d-f3e5bfc51cb4"),
				IntegrationLabel: aws.String("Development AWS Account"),
				IntegrationType:  aws.String("aws-scan"),
				ScanIntervalMins: aws.Int(60),
				ScanEnabled:      aws.Bool(true),
			},
			SourceIntegrationStatus: &snapshotapimodels.SourceIntegrationStatus{
				ScanStatus: aws.String("scanning"),
			},
			SourceIntegrationScanInformation: &snapshotapimodels.SourceIntegrationScanInformation{
				LastScanStartTime: aws.Time(time.Now().Add(time.Duration(-65) * time.Minute)),
				// Last time it scanned was a day ago.
				LastScanEndTime: aws.Time(time.Now().Add(time.Duration(-24) * time.Hour)),
			},
		},
	}
)

func TestPollAndIssueNewScansNoneToRun(t *testing.T) {
	mockLambda := &mockLambdaClient{}

	mockLambda.
		On("Invoke", getTestInvokeInput()).
		// Pass in the first integration, which won't need a new scan.
		Return(getTestInvokeOutput(exampleIntegrations[:1], 200), nil)
	lambdaClient = mockLambda

	result := PollAndIssueNewScans()

	mockLambda.AssertExpectations(t)
	assert.NoError(t, result)
}

func TestPollAndIssueNewScansZeroIntegrations(t *testing.T) {
	mockLambda := &mockLambdaClient{}
	var emptyOutput []*snapshotapimodels.SourceIntegration

	mockLambda.
		On("Invoke", getTestInvokeInput()).
		// Pass in the first integration, which won't need a new scan.
		Return(getTestInvokeOutput(emptyOutput, 200), nil)
	lambdaClient = mockLambda

	result := PollAndIssueNewScans()

	mockLambda.AssertExpectations(t)
	assert.NoError(t, result)
}

func TestScanIntervalElapsed(t *testing.T) {
	assert.True(t, scanIntervalElapsed(&snapshotapimodels.SourceIntegration{
		SourceIntegrationMetadata: &snapshotapimodels.SourceIntegrationMetadata{
			ScanIntervalMins: aws.Int(30),
		},
		SourceIntegrationScanInformation: &snapshotapimodels.SourceIntegrationScanInformation{
			LastScanEndTime: aws.Time(
				time.Now().Add(time.Duration(-60) * time.Minute),
			),
		},
	}))
}

func TestNewScanNotNeeded(t *testing.T) {
	assert.False(t, scanIntervalElapsed(&snapshotapimodels.SourceIntegration{
		SourceIntegrationMetadata: &snapshotapimodels.SourceIntegrationMetadata{
			ScanIntervalMins: aws.Int(30),
		},
		SourceIntegrationScanInformation: &snapshotapimodels.SourceIntegrationScanInformation{
			LastScanEndTime: aws.Time(
				time.Now().Add(time.Duration(-15) * time.Minute),
			),
		},
	}))
}

func TestScanIsNotOngoingScanning(t *testing.T) {
	assert.False(t, scanIsNotOngoing(&snapshotapimodels.SourceIntegration{
		SourceIntegrationStatus: &snapshotapimodels.SourceIntegrationStatus{
			ScanStatus: aws.String(snapshotapimodels.StatusScanning),
		},
	}))
}

func TestScanIsNotOngoingOK(t *testing.T) {
	assert.True(t, scanIsNotOngoing(&snapshotapimodels.SourceIntegration{
		SourceIntegrationStatus: &snapshotapimodels.SourceIntegrationStatus{
			ScanStatus: aws.String(snapshotapimodels.StatusOK),
		},
	}))
}

func TestGetEnabledIntegrations(t *testing.T) {
	mockLambda := new(mockLambdaClient)
	lambdaClient = mockLambda

	mockLambda.
		On("Invoke", getTestInvokeInput()).
		Return(getTestInvokeOutput(exampleIntegrations, 200), nil)

	integrations, err := getEnabledIntegrations()

	mockLambda.AssertExpectations(t)
	require.NoError(t, err)
	assert.Len(t, integrations, len(exampleIntegrations))
}

func TestGetEnabledIntegrationsError(t *testing.T) {
	mockLambda := new(mockLambdaClient)
	lambdaClient = mockLambda

	mockLambda.
		On("Invoke", getTestInvokeInput()).
		Return(&lambda.InvokeOutput{}, errors.New("fake error"))

	_, err := getEnabledIntegrations()

	mockLambda.AssertExpectations(t)
	require.Error(t, err)
}
