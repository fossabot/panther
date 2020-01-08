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

import "github.com/aws/aws-sdk-go/service/cloudtrail"

const (
	CloudTrailSchema     = "AWS.CloudTrail"
	CloudTrailMetaSchema = "AWS.CloudTrail.Meta"
)

// CloudTrail contains all information about a configured CloudTrail.
//
// This includes the trail info, status, event selectors, and attributes of the logging S3 bucket.
type CloudTrail struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from cloudtrail.Trail
	CloudWatchLogsLogGroupArn  *string
	CloudWatchLogsRoleArn      *string
	HasCustomEventSelectors    *bool
	HomeRegion                 *string
	IncludeGlobalServiceEvents *bool
	IsMultiRegionTrail         *bool
	IsOrganizationTrail        *bool
	KmsKeyId                   *string
	LogFileValidationEnabled   *bool
	S3BucketName               *string
	S3KeyPrefix                *string
	SnsTopicARN                *string
	SnsTopicName               *string // Deprecated by AWS

	// Additional fields
	EventSelectors []*cloudtrail.EventSelector
	Status         *cloudtrail.GetTrailStatusOutput
}

type CloudTrailMeta struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	Trails               []*string
	GlobalEventSelectors []*cloudtrail.EventSelector
}

// CloudTrails are a mapping of all Trails in an account keyed by ARN.
type CloudTrails map[string]*CloudTrail
