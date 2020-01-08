package models

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
