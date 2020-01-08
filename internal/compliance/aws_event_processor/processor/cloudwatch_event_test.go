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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

// drop event if the source is not supported
func TestClassifyCloudWatchEventBadSource(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	require.Nil(t, classifyCloudTrailLog(`{"eventSource": "aws.nuka", "eventType": "AwsApiCall"}`))

	expected := []observer.LoggedEntry{
		{
			Entry:   zapcore.Entry{Level: zapcore.DebugLevel, Message: "dropping event from unsupported source"},
			Context: []zapcore.Field{zap.String("eventSource", "aws.nuka")},
		},
	}
	assert.Equal(t, expected, logs.AllUntimed())
}

// drop event if it describes a failed API call
func TestClassifyCloudWatchEventErrorCode(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	require.Nil(t, classifyCloudTrailLog(
		`{"detail": {"errorCode": "AccessDeniedException", "eventSource": "s3.amazonaws.com"}, "detail-type": "AWS API Call via CloudTrail"}`),
	)

	expected := []observer.LoggedEntry{
		{
			Entry: zapcore.Entry{Level: zapcore.DebugLevel, Message: "dropping failed event"},
			Context: []zapcore.Field{
				zap.String("eventSource", "s3.amazonaws.com"),
				zap.String("errorCode", "AccessDeniedException"),
			},
		},
	}
	assert.Equal(t, expected, logs.AllUntimed())
}

// drop event if its read-only
func TestClassifyCloudWatchEventReadOnly(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	require.Nil(t, classifyCloudTrailLog(
		`{"detail": {"eventName": "ListBuckets", "eventSource": "s3.amazonaws.com"}, "detail-type": "AWS API Call via CloudTrail"}`),
	)

	expected := []observer.LoggedEntry{
		{
			Entry:   zapcore.Entry{Level: zapcore.DebugLevel, Message: "s3.amazonaws.com: ignoring read-only event"},
			Context: []zapcore.Field{zap.String("eventName", "ListBuckets")},
		},
	}
	assert.Equal(t, expected, logs.AllUntimed())
}

// drop event if the service classifier doesn't understand it
func TestClassifyCloudWatchEventClassifyError(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	body :=
		`{	"detail": {
				"eventName": "DeleteBucket",
				"recipientAccountId": "111111111111",
				"eventSource":"s3.amazonaws.com"
			}, 
			"detail-type": "AWS API Call via CloudTrail"}`
	require.Nil(t, classifyCloudTrailLog(body))

	expected := []observer.LoggedEntry{
		{
			Entry: zapcore.Entry{Level: zapcore.ErrorLevel, Message: "s3: empty bucket name"},
			Context: []zapcore.Field{
				zap.String("eventName", "DeleteBucket"),
			},
		},
	}
	assert.Equal(t, expected, logs.AllUntimed())
}

// drop event if the account ID is not recognized
func TestClassifyCloudWatchEventUnauthorized(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	body := `{"eventType" : "AwsApiCall", "eventSource": "s3.amazonaws.com", "requestParameters": {"bucketName": "panther"}}`
	changes := classifyCloudTrailLog(body)
	assert.Len(t, changes, 0)

	expected := []observer.LoggedEntry{
		{
			Entry: zapcore.Entry{Level: zapcore.WarnLevel, Message: "dropping event from unauthorized account"},
			Context: []zapcore.Field{
				zap.String("accountId", ""),
				zap.String("eventSource", "s3.amazonaws.com"),
			},
		},
	}
	assert.Equal(t, expected, logs.AllUntimed())
}

func TestClassifyCloudWatchEvent(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	body := `
{
	"detail-type": "AWS API Call via CloudTrail"
    "detail": {
		"recipientAccountId": "111111111111",
    	"eventSource": "s3.amazonaws.com",
        "awsRegion": "us-west-2",
        "eventName": "DeleteBucket",
        "eventTime": "2019-08-01T04:43:00Z",
        "requestParameters": {"bucketName": "panther"},
		"userIdentity": {"accountId": "111111111111"}
    }
}`
	result := classifyCloudTrailLog(body)
	expected := []*resourceChange{{
		AwsAccountID:  "111111111111",
		Delete:        true,
		EventName:     "DeleteBucket",
		EventTime:     "2019-08-01T04:43:00Z",
		IntegrationID: "ebb4d69f-177b-4eff-a7a6-9251fdc72d21",
		ResourceID:    "arn:aws:s3:::panther",
		ResourceType:  schemas.S3BucketSchema,
	}}
	assert.Equal(t, expected, result)
	assert.Empty(t, logs.AllUntimed())
}
