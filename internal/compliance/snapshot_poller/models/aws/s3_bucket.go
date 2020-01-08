package aws

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

import "github.com/aws/aws-sdk-go/service/s3"

// S3BucketSchema is the name of the S3Bucket Schema
const S3BucketSchema = "AWS.S3.Bucket"

// S3Bucket contains all information about an S3 bucket.
type S3Bucket struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	EncryptionRules                []*s3.ServerSideEncryptionRule
	Grants                         []*s3.Grant
	LifecycleRules                 []*s3.LifecycleRule
	LoggingPolicy                  *s3.LoggingEnabled
	MFADelete                      *string
	ObjectLockConfiguration        *s3.ObjectLockConfiguration
	Owner                          *s3.Owner
	Policy                         *string
	PublicAccessBlockConfiguration *s3.PublicAccessBlockConfiguration
	Versioning                     *string
}
