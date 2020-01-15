package api

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
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/panther-labs/panther/api/lambda/snapshot/models"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb"
	"github.com/panther-labs/panther/internal/compliance/snapshot_api/ddb/modelstest"
	pollermodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/poller"
	awspoller "github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws"
)

// Mocks

// mockSQSClient mocks API calls to SQS.
type mockSQSClient struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (client *mockSQSClient) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	args := client.Called(input)
	return args.Get(0).(*sqs.SendMessageBatchOutput), args.Error(1)
}

func (client *mockSQSClient) AddPermission(input *sqs.AddPermissionInput) (*sqs.AddPermissionOutput, error) {
	args := client.Called(input)
	return args.Get(0).(*sqs.AddPermissionOutput), args.Error(1)
}

func (client *mockSQSClient) RemovePermission(input *sqs.RemovePermissionInput) (*sqs.RemovePermissionOutput, error) {
	args := client.Called(input)
	return args.Get(0).(*sqs.RemovePermissionOutput), args.Error(1)
}

func generateMockSQSBatchInputOutput(integrations []*models.SourceIntegrationMetadata) (
	*sqs.SendMessageBatchInput, *sqs.SendMessageBatchOutput, error) {

	// Setup input/output
	var sqsEntries []*sqs.SendMessageBatchRequestEntry
	var err error
	in := &sqs.SendMessageBatchInput{
		QueueUrl: aws.String("test-url"),
	}
	out := &sqs.SendMessageBatchOutput{
		Successful: []*sqs.SendMessageBatchResultEntry{
			{
				Id:               integrations[0].IntegrationID,
				MessageId:        integrations[0].IntegrationID,
				MD5OfMessageBody: aws.String("f6255bb01c648fe967714d52a89e8e9c"),
			},
		},
	}

	// Generate all messages for scans
	for _, integration := range integrations {
		for resourceType := range awspoller.ServicePollers {
			scanMsg := &pollermodels.ScanMsg{
				Entries: []*pollermodels.ScanEntry{
					{
						AWSAccountID:  integration.AWSAccountID,
						IntegrationID: integration.IntegrationID,
						ResourceType:  aws.String(resourceType),
					},
				},
			}

			var messageBodyBytes []byte
			messageBodyBytes, err = jsoniter.Marshal(scanMsg)
			if err != nil {
				break
			}

			sqsEntries = append(sqsEntries, &sqs.SendMessageBatchRequestEntry{
				Id:          integration.IntegrationID,
				MessageBody: aws.String(string(messageBodyBytes)),
			})
		}
	}

	in.Entries = sqsEntries
	return in, out, err
}

// Unit Tests

func TestAddToSnapshotQueue(t *testing.T) {
	snapshotPollersQueueURL = "test-url"
	testIntegrations := []*models.SourceIntegrationMetadata{
		{
			AWSAccountID:     aws.String(testAccountID),
			CreatedAtTime:    aws.Time(time.Time{}),
			CreatedBy:        aws.String("Bobert"),
			IntegrationID:    aws.String(testIntegrationID),
			IntegrationLabel: aws.String("BobertTest"),
			IntegrationType:  aws.String("aws-scan"),
			ScanEnabled:      aws.Bool(true),
			ScanIntervalMins: aws.Int(60),
		},
	}

	sqsIn, sqsOut, err := generateMockSQSBatchInputOutput(testIntegrations)
	require.NoError(t, err)

	mockSQS := &mockSQSClient{}
	// It's non trivial to mock when the order of a slice is not promised
	mockSQS.On("SendMessageBatch", mock.Anything).Return(sqsOut, nil)
	SQSClient = mockSQS

	err = ScanAllResources(testIntegrations)

	require.NoError(t, err)
	// Check that there is one message per service
	assert.Len(t, sqsIn.Entries, len(awspoller.ServicePollers))
}

func TestPutIntegration(t *testing.T) {
	mockSQS := &mockSQSClient{}
	mockSQS.On("SendMessageBatch", mock.Anything).Return(&sqs.SendMessageBatchOutput{}, nil)
	SQSClient = mockSQS
	db = &ddb.DDB{Client: &modelstest.MockDDBClient{TestErr: false}, TableName: "test"}

	out, err := apiTest.PutIntegration(&models.PutIntegrationInput{
		Integrations: []*models.PutIntegrationSettings{
			{
				AWSAccountID:     aws.String(testAccountID),
				IntegrationLabel: aws.String(testIntegrationLabel),
				IntegrationType:  aws.String(testIntegrationType),
				ScanEnabled:      aws.Bool(true),
				ScanIntervalMins: aws.Int(60),
				UserID:           aws.String(testUserID),
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func TestPutIntegrationValidInput(t *testing.T) {
	validator, err := models.Validator()
	require.NoError(t, err)
	assert.NoError(t, validator.Struct(&models.PutIntegrationInput{
		Integrations: []*models.PutIntegrationSettings{
			{
				AWSAccountID:     aws.String(testAccountID),
				IntegrationLabel: aws.String(testIntegrationLabel),
				IntegrationType:  aws.String(testIntegrationType),
				ScanEnabled:      aws.Bool(true),
				ScanIntervalMins: aws.Int(60),
				UserID:           aws.String(testUserID),
			},
		},
	}))
}

func TestPutIntegrationInvalidInput(t *testing.T) {
	validator, err := models.Validator()
	require.NoError(t, err)
	assert.Error(t, validator.Struct(&models.PutIntegrationInput{
		Integrations: []*models.PutIntegrationSettings{
			{
				// Long account ID
				AWSAccountID: aws.String("11111111111111"),
				ScanEnabled:  aws.Bool(true),
				// Invalid integration type
				IntegrationType: aws.String("type-that-does-not-exist"),
			},
		},
	}))
}

func TestPutIntegrationDatabaseError(t *testing.T) {
	in := &models.PutIntegrationInput{
		Integrations: []*models.PutIntegrationSettings{
			{
				AWSAccountID:     aws.String(testAccountID),
				IntegrationLabel: aws.String(testIntegrationLabel),
				IntegrationType:  aws.String(testIntegrationType),
				ScanEnabled:      aws.Bool(true),
				UserID:           aws.String(testUserID),
			},
		},
	}

	db = &ddb.DDB{
		Client: &modelstest.MockDDBClient{
			TestErr: true,
		},
		TableName: "test",
	}

	mockSQS := &mockSQSClient{}
	SQSClient = mockSQS
	mockSQS.On("AddPermission", mock.Anything).Return(&sqs.AddPermissionOutput{}, nil)
	// RemoveRermission will be called to remove the permission that was added previously
	// This is done as part of rollback process to bring the system in a consistent state
	mockSQS.On("RemovePermission", mock.Anything).Return(&sqs.RemovePermissionOutput{}, nil)

	out, err := apiTest.PutIntegration(in)
	assert.Error(t, err)
	assert.Empty(t, out)
}

func TestPutIntegrationDatabaseErrorRecoveryFails(t *testing.T) {
	// Used to capture logs for unit testing purposes
	core, recordedLogs := observer.New(zapcore.ErrorLevel)
	zap.ReplaceGlobals(zap.New(core))

	in := &models.PutIntegrationInput{
		Integrations: []*models.PutIntegrationSettings{
			{
				AWSAccountID:     aws.String(testAccountID),
				IntegrationLabel: aws.String(testIntegrationLabel),
				IntegrationType:  aws.String(testIntegrationType),
				ScanEnabled:      aws.Bool(true),
				UserID:           aws.String(testUserID),
			},
		},
	}

	db = &ddb.DDB{
		Client: &modelstest.MockDDBClient{
			TestErr: true,
		},
		TableName: "test",
	}

	mockSQS := &mockSQSClient{}
	SQSClient = mockSQS
	mockSQS.On("AddPermission", mock.Anything).Return(&sqs.AddPermissionOutput{}, nil)
	// RemoveRermission will be called to remove the permission that was added previously
	// This is done as part of rollback process to bring the system in a consistent state
	mockSQS.On("RemovePermission", mock.Anything).Return(&sqs.RemovePermissionOutput{}, errors.New("error"))

	out, err := apiTest.PutIntegration(in)
	require.Error(t, err)
	require.Empty(t, out)

	errorLog := recordedLogs.FilterMessage("failed to remove SQS permission for integration." +
		" SQS queue has additional permissions that have to be removed manually")
	require.NotNil(t, errorLog)
}

func TestPutLogIntegrationUpdateSqsQueuePermissions(t *testing.T) {
	mockSQS := &mockSQSClient{}
	SQSClient = mockSQS
	logAnalysisQueueURL = "https://sqs.eu-west-1.amazonaws.com/123456789012/testqueue"

	mockSQS.On("AddPermission", mock.Anything).Return(&sqs.AddPermissionOutput{}, nil)
	db = &ddb.DDB{Client: &modelstest.MockDDBClient{TestErr: false}, TableName: "test"}

	out, err := apiTest.PutIntegration(&models.PutIntegrationInput{
		Integrations: []*models.PutIntegrationSettings{
			{
				AWSAccountID:    aws.String(testAccountID),
				IntegrationType: aws.String(models.IntegrationTypeAWS3),
				UserID:          aws.String(testUserID),
				S3Buckets:       aws.StringSlice([]string{"bucket"}),
				KmsKeys:         aws.StringSlice([]string{"keyarns"}),
			},
		},
	})

	sqsPermissionInput := mockSQS.Calls[0].Arguments[0].(*sqs.AddPermissionInput)
	require.Equal(t, aws.StringSlice([]string{testAccountID}), sqsPermissionInput.AWSAccountIds)
	require.Equal(t, aws.StringSlice([]string{"ReceiveMessage"}), sqsPermissionInput.Actions)
	require.Equal(t, aws.String(logAnalysisQueueURL), sqsPermissionInput.QueueUrl)
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func TestPutLogIntegrationUpdateSqsQueuePermissionsFailure(t *testing.T) {
	mockSQS := &mockSQSClient{}
	SQSClient = mockSQS
	logAnalysisQueueURL = "https://sqs.eu-west-1.amazonaws.com/123456789012/testqueue"

	mockSQS.On("AddPermission", mock.Anything).Return(&sqs.AddPermissionOutput{}, errors.New("error"))
	db = &ddb.DDB{Client: &modelstest.MockDDBClient{TestErr: false}, TableName: "test"}

	out, err := apiTest.PutIntegration(&models.PutIntegrationInput{
		Integrations: []*models.PutIntegrationSettings{
			{
				AWSAccountID:    aws.String(testAccountID),
				IntegrationType: aws.String(models.IntegrationTypeAWS3),
				UserID:          aws.String(testUserID),
				S3Buckets:       aws.StringSlice([]string{"bucket"}),
				KmsKeys:         aws.StringSlice([]string{"keyarns"}),
			},
		},
	})
	require.Error(t, err)
	require.Empty(t, out)
}
