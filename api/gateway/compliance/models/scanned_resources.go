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
	"strconv"

	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ScannedResources scanned resources
// swagger:model ScannedResources
type ScannedResources struct {

	// by type
	// Required: true
	ByType []*ResourceOfType `json:"byType"`
}

// Validate validates this scanned resources
func (m *ScannedResources) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateByType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ScannedResources) validateByType(formats strfmt.Registry) error {

	if err := validate.Required("byType", "body", m.ByType); err != nil {
		return err
	}

	for i := 0; i < len(m.ByType); i++ {
		if swag.IsZero(m.ByType[i]) { // not required
			continue
		}

		if m.ByType[i] != nil {
			if err := m.ByType[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("byType" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ScannedResources) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ScannedResources) UnmarshalBinary(b []byte) error {
	var res ScannedResources
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
