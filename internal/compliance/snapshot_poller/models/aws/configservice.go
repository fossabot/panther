package aws

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
