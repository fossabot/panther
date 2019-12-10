package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	schemas "github.com/panther-labs/panther/api/snapshot/aws"
)

// drop event if the source is not supported
func TestClassifyCloudWatchEventBadSource(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	require.Nil(t, classifyCloudWatchEvent(`{"source": "aws.nuka"}`))

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
	require.Nil(t, classifyCloudWatchEvent(`{"source": "aws.s3", "detail": {"errorCode": "AccessDeniedException"}}`))

	expected := []observer.LoggedEntry{
		{
			Entry: zapcore.Entry{Level: zapcore.DebugLevel, Message: "dropping failed event"},
			Context: []zapcore.Field{
				zap.String("eventSource", "aws.s3"),
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
	require.Nil(t, classifyCloudWatchEvent(`{"source": "aws.s3", "detail": {"eventName": "ListBuckets"}`))

	expected := []observer.LoggedEntry{
		{
			Entry:   zapcore.Entry{Level: zapcore.DebugLevel, Message: "aws.s3: ignoring read-only event"},
			Context: []zapcore.Field{zap.String("eventName", "ListBuckets")},
		},
	}
	assert.Equal(t, expected, logs.AllUntimed())
}

// drop event if the service classifier doesn't understand it
func TestClassifyCloudWatchEventClassifyError(t *testing.T) {
	logs := mockLogger()
	accounts = exampleAccounts
	body := `{"source": "aws.s3", "account": "111111111111", "detail": {"eventName": "DeleteBucket"}}`
	require.Nil(t, classifyCloudWatchEvent(body))

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
	body := `{"source": "aws.s3", "detail": {"requestParameters": {"bucketName": "panther"}}}`
	changes := classifyCloudWatchEvent(body)
	assert.Len(t, changes, 0)

	expected := []observer.LoggedEntry{
		{
			Entry: zapcore.Entry{Level: zapcore.WarnLevel, Message: "dropping event from unauthorized account"},
			Context: []zapcore.Field{
				zap.String("accountId", ""),
				zap.String("eventSource", "aws.s3"),
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
    "source": "aws.s3",
	"account": "111111111111"
    "detail": {
        "awsRegion": "us-west-2",
        "eventName": "DeleteBucket",
        "eventTime": "2019-08-01T04:43:00Z",
        "requestParameters": {"bucketName": "panther"},
		"userIdentity": {"accountId": "111111111111"}
    }
}`
	result := classifyCloudWatchEvent(body)
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
