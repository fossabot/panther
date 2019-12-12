package aws

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/go-openapi/strfmt"

	resourcesapimodels "github.com/panther-labs/panther/api/resources/models"
)

// Used to populate the GenericAWSResource.Region field for global AWS resources
const GlobalRegion = "global"

// GenericResource contains fields that will be common to all resources, at some point this will
// probably exist in a more global package but for now since this is the only poller it will exist
// here.
type GenericResource struct {
	ResourceID   *string          `json:"ResourceId"`   // A panther wide unique identifier
	ResourceType *string          `json:"ResourceType"` // A panther defined resource type
	TimeCreated  *strfmt.DateTime `json:"TimeCreated"`  // A standardized format for when the resource was created
}

// GenericAWSResource contains information that is standard across AWS resources
type GenericAWSResource struct {
	//
	// The fields ARN, ID, and Name are tagged omitempty so that those fields will not exist in the
	// down stream python resources if they are not populated, as those fields will either always
	// exist or never exist for a given resource type.
	//
	// This is in contrast to the fields AccountID and Region which will always exist and should therefore
	// always be sent downstream, and the field Tags which may or may not exist for a specific
	// resource of a given resource type and should therefore always be sent downstream, even if
	// only as an empty list.
	//

	// Fields that generally need to be populated after building the snapshot
	AccountID *string `json:"AccountId"` // The ID of the AWS Account the resource resides in
	Region    *string `json:"Region"`    // The region the resource exists in, value of GLOBAL_REGION if global

	// Fields that can generally be populated while building the snapshot
	ARN  *string            `json:"Arn,omitempty"`  // The Amazon Resource Name (ARN)
	ID   *string            `json:"Id,omitempty"`   // The AWS resource identifier
	Name *string            `json:"Name,omitempty"` // The AWS resource name
	Tags map[string]*string // A standardized format for key/value resource tags
}

// ResourcePollerInput contains the metadata to request AWS resource info.
type ResourcePollerInput struct {
	AuthSource          *string
	AuthSourceParsedARN arn.ARN
	IntegrationID       *string
	Regions             []*string
	Timestamp           *strfmt.DateTime
}

// ResourcePoller represents a function to poll a specific AWS resource.
type ResourcePoller func(input *ResourcePollerInput) ([]*resourcesapimodels.AddResourceEntry, error)
