// Code generated by go-swagger; DO NOT EDIT.

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

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RemediateResource remediate resource
// swagger:model RemediateResource
type RemediateResource struct {

	// policy Id
	// Required: true
	PolicyID PolicyID `json:"policyId"`

	// resource Id
	// Required: true
	ResourceID ResourceID `json:"resourceId"`
}

// Validate validates this remediate resource
func (m *RemediateResource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePolicyID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RemediateResource) validatePolicyID(formats strfmt.Registry) error {

	if err := m.PolicyID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("policyId")
		}
		return err
	}

	return nil
}

func (m *RemediateResource) validateResourceID(formats strfmt.Registry) error {

	if err := m.ResourceID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("resourceId")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RemediateResource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RemediateResource) UnmarshalBinary(b []byte) error {
	var res RemediateResource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
