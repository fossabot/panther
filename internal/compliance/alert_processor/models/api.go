package models

import (
	"time"
)

// ComplianceNotification represents the event sent to the AlertProcessor by the compliance engine.
type ComplianceNotification struct {

	//ResourceID is the ID specific to the resource
	ResourceID *string `json:"resourceId" validate:"required,min=1"`

	//PolicyID is the id of the policy that triggered
	PolicyID *string `json:"policyId" validate:"required,min=1"`

	//PolicyVersionID is the version of policy when the alert triggered
	PolicyVersionID *string `json:"policyVersionId"`

	//ShouldAlert indicates whether this notification should cause an alert to be send to the customer
	ShouldAlert *bool `json:"shouldAlert"`

	//Timestamp indicates when the policy was actually evaluated
	Timestamp *time.Time `json:"timestamp"`
}
