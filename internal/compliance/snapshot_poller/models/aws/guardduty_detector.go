package aws

import (
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/go-openapi/strfmt"
)

const (
	GuardDutySchema     = "AWS.GuardDuty.Detector"
	GuardDutyMetaSchema = "AWS.GuardDuty.Detector.Meta"
)

// GuardDutyDetector contains information about a GuardDuty Detector
type GuardDutyDetector struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from guardduty.GetDetectorOutput
	FindingPublishingFrequency *string
	ServiceRole                *string
	Status                     *string
	UpdatedAt                  *strfmt.DateTime

	// Additional fields
	Master *guardduty.Master
}

// GuardDutyMeta contains metadata about all GuardDuty detectors in an account.
type GuardDutyMeta struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Additional fields
	Detectors []*string
}
