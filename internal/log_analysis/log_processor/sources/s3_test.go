package sources

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

func TestParseCloudTrailNotification(t *testing.T) {
	notification := "{\"s3Bucket\": \"testbucket\", \"s3ObjectKey\": [\"key1\",\"key2\"]}"
	expectedOutput := []*S3ObjectInfo{
		{
			S3Bucket:    aws.String("testbucket"),
			S3ObjectKey: aws.String("key1"),
		},
		{
			S3Bucket:    aws.String("testbucket"),
			S3ObjectKey: aws.String("key2"),
		},
	}
	s3Objects, err := ParseNotification(notification)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, s3Objects)
}

func TestParseS3Notification(t *testing.T) {
	//nolint:lll
	notification := "{\"Records\":[{\"eventVersion\":\"2.1\",\"eventSource\":\"aws:s3\",\"awsRegion\":\"us-west-2\",\"eventTime\":\"1970-01-01T00:00:00.000Z\"," +
		"\"eventName\":\"ObjectCreated:Put\",\"userIdentity\":{\"principalId\":\"AIDAJDPLRKLG7UEXAMPLE\"},\"requestParameters\":{\"sourceIPAddress\":\"127.0.0.1\"}," +
		"\"responseElements\":{\"x-amz-request-id\":\"C3D13FE58DE4C810\",\"x-amz-id-2\":\"FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD\"}," +
		"\"s3\":{\"s3SchemaVersion\":\"1.0\",\"configurationId\":\"testConfigRule\"," +
		"\"bucket\":{\"name\":\"mybucket\",\"ownerIdentity\":{\"principalId\":\"A3NL1KOZZKExample\"},\"arn\":\"arn:aws:s3:::mybucket\"},\"object\":{\"key\":\"key1\",\"size\":1024," +
		"\"eTag\":\"d41d8cd98f00b204e9800998ecf8427e\",\"versionId\":\"096fKKXTRTtl3on89fVO.nfljtsv6qko\",\"sequencer\":\"0055AED6DCD90281E5\"}}}]}"
	expectedOutput := []*S3ObjectInfo{
		{
			S3Bucket:    aws.String("mybucket"),
			S3ObjectKey: aws.String("key1"),
		},
	}
	s3Objects, err := ParseNotification(notification)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, s3Objects)
}
