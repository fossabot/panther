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
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	schemas "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
)

func classifyS3(detail gjson.Result, accountID string) []*resourceChange {
	eventName := detail.Get("eventName").Str

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/list_amazons3.html
	if eventName == "UploadPart" ||
		eventName == "CreateMultipartUpload" ||
		eventName == "CompleteMultipartUpload" ||
		eventName == "HeadBucket" ||
		eventName == "PutObject" {

		zap.L().Debug("s3: ignoring event", zap.String("eventName", eventName))
		return nil
	}

	bucketName := detail.Get("requestParameters.bucketName").Str
	if bucketName == "" {
		zap.L().Error("s3: empty bucket name", zap.String("eventName", eventName))
		return nil
	}

	s3ARN := arn.ARN{
		Partition: "aws",
		Service:   "s3",
		Region:    "",
		AccountID: "",
		Resource:  bucketName,
	}

	return []*resourceChange{{
		// Incredibly, CloudTrail logs do not indicate which account owns the bucket.
		//
		// We would fall back to the account which generated the log, but the "recipientAccountId"
		// normally in CloudTrail logs doesn't appear to be populated in CloudWatch events.
		//
		// The *only* place in the log that contains an account number is the user identity.
		// So we assume that the user or role making the API call lives in the same account as the
		// the bucket itself, which is usually true.
		// TODO - test and document possible exceptions, or pull accountID from SNS wrapper
		//
		// If we are wrong, either the poller fails to describe the bucket
		// (and gives up eventually), or we show a bucket as if its part of their account
		// (which in a sense it kind of is - they have read and write access to it).
		AwsAccountID: accountID,

		Delete:    eventName == "DeleteBucket",
		EventName: eventName,
		// Format: arn:aws:s3:::bucket_name
		ResourceID:   s3ARN.String(),
		ResourceType: schemas.S3BucketSchema,
	}}
}
