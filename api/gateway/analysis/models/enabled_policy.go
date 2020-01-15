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

// EnabledPolicy enabled policy
// swagger:model EnabledPolicy
type EnabledPolicy struct {

	// body
	Body Body `json:"body,omitempty"`

	// id
	ID ID `json:"id,omitempty"`

	// resource types
	ResourceTypes TypeSet `json:"resourceTypes,omitempty"`

	// severity
	Severity Severity `json:"severity,omitempty"`

	// suppressions
	Suppressions Suppressions `json:"suppressions,omitempty"`

	// version Id
	VersionID VersionID `json:"versionId,omitempty"`
}

// Validate validates this enabled policy
func (m *EnabledPolicy) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBody(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceTypes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSeverity(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSuppressions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersionID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EnabledPolicy) validateBody(formats strfmt.Registry) error {

	if swag.IsZero(m.Body) { // not required
		return nil
	}

	if err := m.Body.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("body")
		}
		return err
	}

	return nil
}

func (m *EnabledPolicy) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := m.ID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("id")
		}
		return err
	}

	return nil
}

func (m *EnabledPolicy) validateResourceTypes(formats strfmt.Registry) error {

	if swag.IsZero(m.ResourceTypes) { // not required
		return nil
	}

	if err := m.ResourceTypes.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("resourceTypes")
		}
		return err
	}

	return nil
}

func (m *EnabledPolicy) validateSeverity(formats strfmt.Registry) error {

	if swag.IsZero(m.Severity) { // not required
		return nil
	}

	if err := m.Severity.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("severity")
		}
		return err
	}

	return nil
}

func (m *EnabledPolicy) validateSuppressions(formats strfmt.Registry) error {

	if swag.IsZero(m.Suppressions) { // not required
		return nil
	}

	if err := m.Suppressions.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("suppressions")
		}
		return err
	}

	return nil
}

func (m *EnabledPolicy) validateVersionID(formats strfmt.Registry) error {

	if swag.IsZero(m.VersionID) { // not required
		return nil
	}

	if err := m.VersionID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("versionId")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *EnabledPolicy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EnabledPolicy) UnmarshalBinary(b []byte) error {
	var res EnabledPolicy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
