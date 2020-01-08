package awslogs

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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

func TestCloudTrailLog(t *testing.T) {
	parser := &CloudTrailParser{}

	//nolint:lll
	log := `{"Records": [{"eventVersion":"1.05","userIdentity":{"type":"AWSService","invokedBy":"cloudtrail.amazonaws.com"},"eventTime":"2018-08-26T14:17:23Z","eventSource":"kms.amazonaws.com","eventName":"GenerateDataKey","awsRegion":"us-west-2","sourceIPAddress":"cloudtrail.amazonaws.com","userAgent":"cloudtrail.amazonaws.com","requestParameters":{"keySpec":"AES_256","encryptionContext":{"aws:cloudtrail:arn":"arn:aws:cloudtrail:us-west-2:888888888888:trail/panther-lab-cloudtrail","aws:s3:arn":"arn:aws:s3:::panther-lab-cloudtrail/AWSLogs/888888888888/CloudTrail/us-west-2/2018/08/26/888888888888_CloudTrail_us-west-2_20180826T1410Z_inUwlhwpSGtlqmIN.json.gz"},"keyId":"arn:aws:kms:us-west-2:888888888888:key/72c37aae-1000-4058-93d4-86374c0fe9a0"},"responseElements":null,"requestID":"3cff2472-5a91-4bd9-b6d2-8a7a1aaa9086","eventID":"7a215e16-e0ad-4f6c-82b9-33ff6bbdedd2","readOnly":true,"resources":[{"ARN":"arn:aws:kms:us-west-2:888888888888:key/72c37aae-1000-4058-93d4-86374c0fe9a0","accountId":"888888888888","type":"AWS::KMS::Key"}],"eventType":"AwsApiCall","recipientAccountId":"888888888888","sharedEventID":"238c190c-1a30-4756-8e08-19fc36ad1b9f"}]}`

	expectedDate := time.Unix(1535293043, 0).In(time.UTC)
	expectedEvent := &CloudTrail{
		EventVersion: aws.String("1.05"),
		UserIdentity: &CloudTrailUserIdentity{
			Type:      aws.String("AWSService"),
			InvokedBy: aws.String("cloudtrail.amazonaws.com"),
		},
		EventTime:       (*timestamp.RFC3339)(&expectedDate),
		EventSource:     aws.String("kms.amazonaws.com"),
		EventName:       aws.String("GenerateDataKey"),
		AWSRegion:       aws.String("us-west-2"),
		SourceIPAddress: aws.String("cloudtrail.amazonaws.com"),
		UserAgent:       aws.String("cloudtrail.amazonaws.com"),
		RequestID:       aws.String("3cff2472-5a91-4bd9-b6d2-8a7a1aaa9086"),
		EventID:         aws.String("7a215e16-e0ad-4f6c-82b9-33ff6bbdedd2"),
		ReadOnly:        aws.Bool(true),
		Resources: []CloudTrailResources{
			{
				ARN:       aws.String("arn:aws:kms:us-west-2:888888888888:key/72c37aae-1000-4058-93d4-86374c0fe9a0"),
				AccountID: aws.String("888888888888"),
				Type:      aws.String("AWS::KMS::Key"),
			},
		},
		EventType:          aws.String("AwsApiCall"),
		RecipientAccountID: aws.String("888888888888"),
		SharedEventID:      aws.String("238c190c-1a30-4756-8e08-19fc36ad1b9f"),
		RequestParameters: map[string]interface{}{
			"keyId":   "arn:aws:kms:us-west-2:888888888888:key/72c37aae-1000-4058-93d4-86374c0fe9a0",
			"keySpec": "AES_256",
			"encryptionContext": map[string]interface{}{
				"aws:cloudtrail:arn": "arn:aws:cloudtrail:us-west-2:888888888888:trail/panther-lab-cloudtrail",
				//nolint:lll
				"aws:s3:arn": "arn:aws:s3:::panther-lab-cloudtrail/AWSLogs/888888888888/CloudTrail/us-west-2/2018/08/26/888888888888_CloudTrail_us-west-2_20180826T1410Z_inUwlhwpSGtlqmIN.json.gz",
			},
		},
	}

	require.Equal(t, (interface{})(expectedEvent), parser.Parse(log)[0])
}

func TestCloudTrailLogType(t *testing.T) {
	parser := &CloudTrailParser{}
	require.Equal(t, "AWS.CloudTrail", parser.LogType())
}
