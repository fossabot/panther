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

import "github.com/aws/aws-sdk-go/service/configservice"

const (
	// ConfigServiceSchema is the schema ID for the ConfigService type.
	ConfigServiceSchema = "AWS.Config.Recorder"
	// ConfigServiceMetaSchema is the schema ID for the ConfigServiceMeta type.
	ConfigServiceMetaSchema = "AWS.Config.Recorder.Meta"
)

// ConfigService contains all information about a policy.
type ConfigService struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from configservice.ConfigurationRecorder
	RecordingGroup *configservice.RecordingGroup
	RoleARN        *string

	// Additional fields
	Status *configservice.ConfigurationRecorderStatus
}

// ConfigServiceMeta contains metadata about all Config Service Recorders in an account.
type ConfigServiceMeta struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	GlobalRecorderCount *int
	Recorders           []*string
}
