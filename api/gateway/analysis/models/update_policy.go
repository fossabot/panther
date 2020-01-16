// Code generated by go-swagger; DO NOT EDIT.

package models

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UpdatePolicy update policy
// swagger:model UpdatePolicy
type UpdatePolicy struct {

	// auto remediation Id
	AutoRemediationID AutoRemediationID `json:"autoRemediationId,omitempty"`

	// auto remediation parameters
	AutoRemediationParameters AutoRemediationParameters `json:"autoRemediationParameters,omitempty"`

	// body
	// Required: true
	Body Body `json:"body"`

	// description
	Description Description `json:"description,omitempty"`

	// display name
	DisplayName DisplayName `json:"displayName,omitempty"`

	// enabled
	// Required: true
	Enabled Enabled `json:"enabled"`

	// id
	// Required: true
	ID ID `json:"id"`

	// reference
	Reference Reference `json:"reference,omitempty"`

	// resource types
	ResourceTypes TypeSet `json:"resourceTypes,omitempty"`

	// runbook
	Runbook Runbook `json:"runbook,omitempty"`

	// severity
	// Required: true
	Severity Severity `json:"severity"`

	// suppressions
	Suppressions Suppressions `json:"suppressions,omitempty"`

	// tags
	Tags Tags `json:"tags,omitempty"`

	// tests
	Tests TestSuite `json:"tests,omitempty"`

	// user Id
	// Required: true
	UserID UserID `json:"userId"`
}

// Validate validates this update policy
func (m *UpdatePolicy) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAutoRemediationID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAutoRemediationParameters(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBody(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDisplayName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnabled(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReference(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceTypes(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRunbook(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSeverity(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSuppressions(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTests(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UpdatePolicy) validateAutoRemediationID(formats strfmt.Registry) error {

	if swag.IsZero(m.AutoRemediationID) { // not required
		return nil
	}

	if err := m.AutoRemediationID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("autoRemediationId")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateAutoRemediationParameters(formats strfmt.Registry) error {

	if swag.IsZero(m.AutoRemediationParameters) { // not required
		return nil
	}

	if err := m.AutoRemediationParameters.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("autoRemediationParameters")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateBody(formats strfmt.Registry) error {

	if err := m.Body.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("body")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateDescription(formats strfmt.Registry) error {

	if swag.IsZero(m.Description) { // not required
		return nil
	}

	if err := m.Description.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("description")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateDisplayName(formats strfmt.Registry) error {

	if swag.IsZero(m.DisplayName) { // not required
		return nil
	}

	if err := m.DisplayName.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("displayName")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateEnabled(formats strfmt.Registry) error {

	if err := m.Enabled.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("enabled")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateID(formats strfmt.Registry) error {

	if err := m.ID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("id")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateReference(formats strfmt.Registry) error {

	if swag.IsZero(m.Reference) { // not required
		return nil
	}

	if err := m.Reference.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("reference")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateResourceTypes(formats strfmt.Registry) error {

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

func (m *UpdatePolicy) validateRunbook(formats strfmt.Registry) error {

	if swag.IsZero(m.Runbook) { // not required
		return nil
	}

	if err := m.Runbook.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("runbook")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateSeverity(formats strfmt.Registry) error {

	if err := m.Severity.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("severity")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateSuppressions(formats strfmt.Registry) error {

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

func (m *UpdatePolicy) validateTags(formats strfmt.Registry) error {

	if swag.IsZero(m.Tags) { // not required
		return nil
	}

	if err := m.Tags.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("tags")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateTests(formats strfmt.Registry) error {

	if swag.IsZero(m.Tests) { // not required
		return nil
	}

	if err := m.Tests.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("tests")
		}
		return err
	}

	return nil
}

func (m *UpdatePolicy) validateUserID(formats strfmt.Registry) error {

	if err := m.UserID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("userId")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UpdatePolicy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdatePolicy) UnmarshalBinary(b []byte) error {
	var res UpdatePolicy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
